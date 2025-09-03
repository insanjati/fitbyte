package repository

import (
	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*model.UserResponse, error) {
	query := `SELECT name, email, preference, weight_unit, height_unit, weight, height, image_uri FROM users WHERE id = $1`

	var user model.UserResponse
	err := r.db.QueryRow(query, id).Scan(
		&user.Name,
		&user.Email,
		&user.Preference,
		&user.WeightUnit,
		&user.HeightUnit,
		&user.Weight,
		&user.Height,
		&user.ImageUri,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
