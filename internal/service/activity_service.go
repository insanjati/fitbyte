package service

import (
	"errors"
	"time"

	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/repository"

	"github.com/google/uuid"
)

// type ActivityService interface {
//     CreateActivity(activityType string, doneAt time.Time, duration int) (*model.Activity, error)
//     FindActivities(limit, offset int, activityType string, doneAtFrom, doneAtTo *time.Time) ([]model.Activity, error)
// }

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

func (s *ActivityService) CreateActivity(userID int, req model.CreateActivityRequest) (*model.Activity, error) {
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

func (s *ActivityService) GetUserActivities(userID int, filter *model.ActivityFilter) ([]model.Activity, error) {
	return s.activityRepo.GetUserActivities(userID, filter)
}
