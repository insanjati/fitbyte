package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/insanjati/fitbyte/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) CheckExistedUserById(id uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = $1`

	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)

	// Log query, id, and result
	fmt.Printf("CheckExistedUserById | Query: %s | ID: %s | Result: %d | Error: %v\n",
		query, id.String(), exists, err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *ActivityRepository) CheckActivityOwnership(userID uuid.UUID, activityID uuid.UUID) (*model.Activity, error) {
	query := `SELECT * FROM activities WHERE id = $1 AND user_id = $2`

	var activity model.Activity
	err := r.db.QueryRow(query, activityID, userID).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.ActivityType,
		&activity.DoneAt,
		&activity.DurationInMinutes,
		&activity.CaloriesBurned,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	// Log query
	fmt.Printf("CheckActivityOwnership | Query: %s | userID: %s | activityID: %s\n", query, userID, activityID)

	if err != nil {
		return nil, err
	}

	return &activity, nil
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

func (r *ActivityRepository) UpdateActivity(userID uuid.UUID, activityID uuid.UUID, updatedAt time.Time, req *model.UpdateActivityRequest, caloriesBurned int) (*model.Activity, error) {
	query := `
		UPDATE activities
		SET updated_at = $1
	`

	args := []interface{}{updatedAt}
	argIndex := 2

	fields := []string{}

	if req.ActivityType != nil {
		fields = append(fields, fmt.Sprintf(" activity_type = $%d", argIndex))
		args = append(args, *req.ActivityType)
		argIndex++
	}

	if req.DoneAt != nil {
		fields = append(fields, fmt.Sprintf(" done_at = $%d", argIndex))
		args = append(args, *req.DoneAt)
		argIndex++
	}

	if req.DurationInMinutes != nil {
		fields = append(fields, fmt.Sprintf(" duration_in_minutes = $%d", argIndex))
		args = append(args, *req.DurationInMinutes)
		argIndex++
	}

	fields = append(fields, fmt.Sprintf(" calories_burned = $%d", argIndex))
	args = append(args, caloriesBurned)
	argIndex++

	// If no fields are set, there's nothing to update (mending di service ndak sih)
	if len(fields) == 0 {
		// nothing to update
		return nil, nil
	}

	// Join fields for SET clause
	if len(fields) > 0 {
		query += ", " + strings.Join(fields, ", ")
	}

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE user_id = $%d AND id = $%d RETURNING id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at", argIndex, argIndex+1)
	// Log query
	fmt.Printf("UpdateActivity | Query: %s | userID: %s | activityID: %s\n", query, userID, activityID)
	args = append(args, userID, activityID)

	// ini nanti ganti QueryRowContext (Get ni sama ndak kek QueryRow?)
	var activity model.Activity
	row := r.db.QueryRow(query, args...)
	err := row.Scan(
		&activity.ID,
		&activity.UserID,
		&activity.ActivityType,
		&activity.DoneAt,
		&activity.DurationInMinutes,
		&activity.CaloriesBurned,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
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
