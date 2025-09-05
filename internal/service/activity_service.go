package service

import (
	"errors"
	"time"

	appErrors "github.com/insanjati/fitbyte/internal/errors"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"


	"github.com/google/uuid"
)

type ActivityService struct {
	activityRepo *repository.ActivityRepository
}

func NewActivityService(activityRepo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepo: activityRepo}
}

func (s *ActivityService) calculateCalories(activityType model.ActivityType, durationInMinutes int) (int, error) {
	calsPerMinute, ok := model.ActivityTypeCalories[activityType]
	if !ok {
		return 0, errors.New("invalid activityType")
	}
	return calsPerMinute * durationInMinutes, nil
}

func (s *ActivityService) CreateActivity(userID uuid.UUID, req model.CreateActivityRequest) (*model.Activity, error) {
	// Validate userID
	isUserExists, err := s.activityRepo.CheckExistedUserById(userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, appErrors.ErrUnauthorized
	}

	// Validate and parse doneAt
	doneAt, err := time.Parse(time.RFC3339, req.DoneAt)
	if err != nil {
		return nil, errors.New("invalid doneAt")
	}

	// Validate type and calculate calories
	calories, err := s.calculateCalories(req.ActivityType, req.DurationInMinutes)
	if err != nil {
		return nil, err
	}

	activity := &model.Activity{
		ID:                uuid.New(),
		ID:                uuid.New(),
		UserID:            userID,
		ActivityType:      req.ActivityType,
		DoneAt:            doneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    calories,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.activityRepo.CreateActivity(activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *ActivityService) GetUserActivities(userID uuid.UUID, filter *model.ActivityFilter) ([]model.Activity, error) {
func (s *ActivityService) GetUserActivities(userID uuid.UUID, filter *model.ActivityFilter) ([]model.Activity, error) {
	return s.activityRepo.GetUserActivities(userID, filter)
}

func (s *ActivityService) UpdateActivity(userID uuid.UUID, activityID uuid.UUID, req model.UpdateActivityRequest) (*model.Activity, error) {
	// Validate userID
	isUserExists, err := s.activityRepo.CheckExistedUserById(userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, appErrors.ErrUnauthorized
	}

	// Validate activity
	existedActivity, err := s.activityRepo.CheckActivityOwnership(userID, activityID)
	if err != nil {
		return nil, err
	}
	if existedActivity == nil {
		return nil, appErrors.ErrForbidden
	}

	// Prepare data update
	if req.ActivityType != nil {
		existedActivity.ActivityType = *req.ActivityType
	}
	if req.DoneAt != nil {
		// Validate and parse doneAt
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

	// Validate type and re-calculate calories
	calories, err := s.calculateCalories(existedActivity.ActivityType, existedActivity.DurationInMinutes)
	if err != nil {
		return nil, err
	}

	// Update activity
	activity, err := s.activityRepo.UpdateActivity(userID, activityID, time.Now(), &req, calories)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) UpdateActivity(userID uuid.UUID, activityID uuid.UUID, req model.UpdateActivityRequest) (*model.Activity, error) {
	// Validate userID
	isUserExists, err := s.activityRepo.CheckExistedUserById(userID)
	if err != nil {
		return nil, err
	}
	if !isUserExists {
		return nil, appErrors.ErrUnauthorized
	}

	// Validate activity
	existedActivity, err := s.activityRepo.CheckActivityOwnership(userID, activityID)
	if err != nil {
		return nil, err
	}
	if existedActivity == nil {
		return nil, appErrors.ErrForbidden
	}

	// Prepare data update
	if req.ActivityType != nil {
		existedActivity.ActivityType = *req.ActivityType
	}
	if req.DoneAt != nil {
		// Validate and parse doneAt
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

	// Validate type and re-calculate calories
	calories, err := s.calculateCalories(existedActivity.ActivityType, existedActivity.DurationInMinutes)
	if err != nil {
		return nil, err
	}

	// Update activity
	activity, err := s.activityRepo.UpdateActivity(userID, activityID, time.Now(), &req, calories)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *ActivityService) DeleteActivity(activityID uuid.UUID, userID uuid.UUID) error {
	return s.activityRepo.DeleteActivity(activityID, userID)
}