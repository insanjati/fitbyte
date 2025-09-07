package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/cache"
	appErrors "github.com/insanjati/fitbyte/internal/errors"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"
)

type ActivityService struct {
	activityRepo *repository.ActivityRepository
	cache        *cache.Redis
}

func NewActivityService(activityRepo *repository.ActivityRepository, cache *cache.Redis) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		cache:        cache,
	}
}

func (s *ActivityService) getUserActivitiesKey(userID uuid.UUID, filter *model.ActivityFilter) string {
	filterHash := ""
	if filter != nil {
		if filter.ActivityType != nil {
			filterHash += fmt.Sprintf("_type_%s", string(*filter.ActivityType))
		}
		if filter.DoneAtFrom != nil {
			filterHash += fmt.Sprintf("_from_%s", *filter.DoneAtFrom)
		}
		if filter.DoneAtTo != nil {
			filterHash += fmt.Sprintf("_to_%s", *filter.DoneAtTo)
		}
		if filter.CaloriesBurnedMin != nil {
			filterHash += fmt.Sprintf("_cal_min_%d", *filter.CaloriesBurnedMin)
		}
		if filter.CaloriesBurnedMax != nil {
			filterHash += fmt.Sprintf("_cal_max_%d", *filter.CaloriesBurnedMax)
		}
		if filter.Limit != nil {
			filterHash += fmt.Sprintf("_limit_%d", *filter.Limit)
		}
		if filter.Offset != nil {
			filterHash += fmt.Sprintf("_offset_%d", *filter.Offset)
		}
	}
	return fmt.Sprintf("user_activities:%s%s", userID.String(), filterHash)
}

func (s *ActivityService) getActivityKey(activityID uuid.UUID) string {
	return fmt.Sprintf("activity:%s", activityID.String())
}

func (s *ActivityService) getUserExistsKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_exists:%s", userID.String())
}

func (s *ActivityService) getUserActivitiesPattern(userID uuid.UUID) string {
	return fmt.Sprintf("user_activities:%s*", userID.String())
}

func (s *ActivityService) calculateCalories(activityType *model.ActivityType, durationInMinutes *int) (*int, error) {
	calsPerMinute, ok := model.ActivityTypeCalories[*activityType]
	if !ok {
		return nil, errors.New("invalid activityType")
	}
	calories := calsPerMinute * (*durationInMinutes)
	return &calories, nil
}

func (s *ActivityService) CreateActivity(ctx context.Context, userID uuid.UUID, req model.CreateActivityRequest) (*model.Activity, error) {
	isUserExists, err := s.checkUserExistsWithCache(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, appErrors.ErrUnauthorized
	}

	doneAt, err := time.Parse(time.RFC3339, req.DoneAt)
	if err != nil {
		return nil, errors.New("invalid doneAt")
	}

	calories, err := s.calculateCalories(&req.ActivityType, &req.DurationInMinutes)
	if err != nil {
		return nil, err
	}

	activity := &model.Activity{
		ID:                uuid.New(),
		UserID:            userID,
		ActivityType:      req.ActivityType,
		DoneAt:            doneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    *calories,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.activityRepo.CreateActivity(activity); err != nil {
		return nil, err
	}

	activityKey := s.getActivityKey(activity.ID)
	_ = s.cache.SetExp(ctx, activityKey, activity, 1*time.Hour)

	pattern := s.getUserActivitiesPattern(userID)
	_ = s.cache.DeletePattern(ctx, pattern)

	return activity, nil
}

func (s *ActivityService) GetUserActivities(ctx context.Context, userID uuid.UUID, filter *model.ActivityFilter) ([]model.Activity, error) {
	cacheKey := s.getUserActivitiesKey(userID, filter)
	var cachedActivities []model.Activity

	err := s.cache.GetAs(ctx, cacheKey, &cachedActivities)
	if err == nil {
		// log.Printf("[CACHE HIT] UserActivities - Key: %s, UserID: %s", cacheKey, userID.String())
		return cachedActivities, nil
	}

	// log.Printf("[CACHE MISS] UserActivities - Key: %s, UserID: %s, Error: %v", cacheKey, userID.String(), err)

	activities, err := s.activityRepo.GetUserActivities(userID, filter)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetExp(ctx, cacheKey, activities, 30*time.Minute)
	// if cacheErr != nil {
	// 	log.Printf("[CACHE SET ERROR] Key: %s, Error: %v", cacheKey, cacheErr)
	// } else {
	// 	log.Printf("[CACHE SET] Key: %s, Records: %d", cacheKey, len(activities))
	// }

	return activities, nil
}

func (s *ActivityService) UpdateActivity(ctx context.Context, userID uuid.UUID, activityID uuid.UUID, req model.UpdateActivityRequest) (*model.Activity, error) {
	isUserExists, err := s.checkUserExistsWithCache(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, appErrors.ErrUnauthorized
	}

	var existedActivity *model.Activity
	activityKey := s.getActivityKey(activityID)

	var cachedActivity model.Activity
	err = s.cache.GetAs(ctx, activityKey, &cachedActivity)
	if err == nil {
		if cachedActivity.UserID != userID {
			return nil, appErrors.ErrForbidden
		}
		existedActivity = &cachedActivity
	} else {
		existedActivity, err = s.activityRepo.CheckActivityOwnership(userID, activityID)
		if err != nil {
			return nil, err
		}
		if existedActivity == nil {
			return nil, appErrors.ErrForbidden
		}
	}

	if req.ActivityType != nil {
		existedActivity.ActivityType = *req.ActivityType
	}
	if req.DoneAt != nil {
		doneAt, err := time.Parse(time.RFC3339, *req.DoneAt)
		if err != nil {
			return nil, errors.New("invalid doneAt")
		}
		existedActivity.DoneAt = doneAt
	}
	if req.DurationInMinutes != nil {
		if *req.DurationInMinutes < 1 {
			return nil, errors.New("durationInMinutes must be >= 1")
		}
		existedActivity.DurationInMinutes = *req.DurationInMinutes
	}

	calories, err := s.calculateCalories(&existedActivity.ActivityType, &existedActivity.DurationInMinutes)
	if err != nil {
		return nil, err
	}

	activity, err := s.activityRepo.UpdateActivity(userID, activityID, time.Now(), &req, *calories)
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetExp(ctx, activityKey, activity, 1*time.Hour)

	pattern := s.getUserActivitiesPattern(userID)
	_ = s.cache.DeletePattern(ctx, pattern)

	return activity, nil
}

func (s *ActivityService) DeleteActivity(ctx context.Context, activityID uuid.UUID, userID uuid.UUID) error {
	err := s.activityRepo.DeleteActivity(activityID, userID)
	if err != nil {
		return err
	}

	activityKey := s.getActivityKey(activityID)
	_ = s.cache.Delete(ctx, activityKey)

	pattern := s.getUserActivitiesPattern(userID)
	_ = s.cache.DeletePattern(ctx, pattern)

	return nil
}

func (s *ActivityService) checkUserExistsWithCache(ctx context.Context, userID uuid.UUID) (bool, error) {
	cacheKey := s.getUserExistsKey(userID)

	var exists bool
	err := s.cache.GetAs(ctx, cacheKey, &exists)
	if err == nil {
		return exists, nil
	}

	isUserExists, err := s.activityRepo.CheckExistedUserById(userID)
	if err != nil {
		return false, err
	}

	_ = s.cache.SetExp(ctx, cacheKey, isUserExists, 5*time.Minute)

	return isUserExists, nil
}

func (s *ActivityService) GetActivityByID(userID uuid.UUID, activityID uuid.UUID) (*model.Activity, error) {
	ctx := context.Background()

	activityKey := s.getActivityKey(activityID)
	var cachedActivity model.Activity

	err := s.cache.GetAs(ctx, activityKey, &cachedActivity)
	if err == nil {
		if cachedActivity.UserID != userID {
			return nil, appErrors.ErrForbidden
		}
		return &cachedActivity, nil
	}

	activity, err := s.activityRepo.CheckActivityOwnership(userID, activityID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, appErrors.ErrForbidden
	}

	_ = s.cache.SetExp(ctx, activityKey, activity, 1*time.Hour)

	return activity, nil
}
