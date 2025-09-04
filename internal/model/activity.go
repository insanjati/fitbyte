package model

import (
	"time"

	"github.com/google/uuid"
)

type ActivityType string

const (
	ActivityTypeWalking    ActivityType = "Walking"
	ActivityTypeYoga       ActivityType = "Yoga"
	ActivityTypeStretching ActivityType = "Stretching"
	ActivityTypeCycling    ActivityType = "Cycling"
	ActivityTypeSwimming   ActivityType = "Swimming"
	ActivityTypeDancing    ActivityType = "Dancing"
	ActivityTypeHiking     ActivityType = "Hiking"
	ActivityTypeRunning    ActivityType = "Running"
	ActivityTypeHIIT       ActivityType = "HIIT"
	ActivityTypeJumpRope   ActivityType = "JumpRope"
)

var ActivityTypeCalories = map[ActivityType]int{
	ActivityTypeWalking:    4,
	ActivityTypeYoga:       4,
	ActivityTypeStretching: 4,
	ActivityTypeCycling:    8,
	ActivityTypeSwimming:   8,
	ActivityTypeDancing:    8,
	ActivityTypeHiking:     10,
	ActivityTypeRunning:    10,
	ActivityTypeHIIT:       10,
	ActivityTypeJumpRope:   10,
}

type Activity struct {
	ID                uuid.UUID    `json:"activityId" db:"id"`
	UserID            int          `json:"userId" db:"user_id"`
	ActivityType      ActivityType `json:"activityType" db:"activity_type"`
	DoneAt            time.Time    `json:"doneAt" db:"done_at"`
	DurationInMinutes int          `json:"durationInMinutes" db:"duration_in_minutes"`
	CaloriesBurned    int          `json:"caloriesBurned" db:"calories_burned"`
	CreatedAt         time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time    `json:"updatedAt" db:"updated_at"`
}

type CreateActivityRequest struct {
	ActivityType      ActivityType `json:"activityType" binding:"required"`
	DoneAt            string       `json:"doneAt" binding:"required"`
	DurationInMinutes int          `json:"durationInMinutes" binding:"required,min=1"`
}

type UpdateActivityRequest struct {
	ActivityType      *ActivityType `json:"activityType"`
	DoneAt            *string       `json:"doneAt"`
	DurationInMinutes *int          `json:"durationInMinutes" binding:"onitemempty,min=1"`
}

type ActivityFilter struct {
	Limit             *int          `form:"limit"`
	Offset            *int          `form:"offset"`
	ActivityType      *ActivityType `form:"activityType"`
	DoneAtFrom        *string       `form:"doneAtFrom"`
	DoneAtTo          *string       `form:"doneAtTo"`
	CaloriesBurnedMin *int          `form:"caloriesBurnedMin"`
	CaloriesBurnedMax *int          `form:"caloriesBurnedMax"`
}
