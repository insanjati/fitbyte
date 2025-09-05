package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/model"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) CreateActivity(activity *model.Activity) error {
	query := `
		INSERT INTO activities (id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(query,
		activity.ID,
		activity.UserID,
		activity.ActivityType,
		activity.DoneAt,
		activity.DurationInMinutes,
		activity.CaloriesBurned,
		activity.CreatedAt,
		activity.UpdatedAt,
	)
	return err
}

func (r *ActivityRepository) GetUserActivities(userID uuid.UUID, filter *model.ActivityFilter) ([]model.Activity, error) {
	query := `
		SELECT id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at
		FROM activities 
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	// Build WHERE clause based on filters
	conditions := []string{}

	// Add optional ActivityType filter
	if filter.ActivityType != nil {
		conditions = append(conditions, fmt.Sprintf("activity_type = $%d", argIndex))
		args = append(args, *filter.ActivityType)
		argIndex++
	}

	// Add optional DoneAtFrom filter
	if filter.DoneAtFrom != nil {
		doneAtFrom, err := time.Parse(time.RFC3339, *filter.DoneAtFrom)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("done_at >= $%d", argIndex))
			args = append(args, doneAtFrom)
			argIndex++
		}
	}

	// Add optional DoneAtTo filter
	if filter.DoneAtTo != nil {
		doneAtTo, err := time.Parse(time.RFC3339, *filter.DoneAtTo)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("done_at <= $%d", argIndex))
			args = append(args, doneAtTo)
			argIndex++
		}
	}

	// Add optional CaloriesBurnedMin filter
	if filter.CaloriesBurnedMin != nil {
		conditions = append(conditions, fmt.Sprintf("calories_burned >= $%d", argIndex))
		args = append(args, *filter.CaloriesBurnedMin)
		argIndex++
	}

	// Add optional CaloriesBurnedMax filter
	if filter.CaloriesBurnedMax != nil {
		conditions = append(conditions, fmt.Sprintf("calories_burned <= $%d", argIndex))
		args = append(args, *filter.CaloriesBurnedMax)
		argIndex++
	}

	// AND
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ORDER BY
	query += " ORDER BY done_at DESC"

	// Add LIMIT and OFFSET
	limit := 5
	offset := 0

	if filter.Limit != nil && *filter.Limit > 0 {
		limit = *filter.Limit
	}

	if filter.Offset != nil && *filter.Offset >= 0 {
		offset = *filter.Offset
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var activities []model.Activity
	err := r.db.Select(&activities, query, args...)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *ActivityRepository) DeleteActivity(activityID uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM activities 
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(query, activityID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("activity not found")
	}

	return nil
}
