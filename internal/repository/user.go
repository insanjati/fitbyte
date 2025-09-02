package repository

import (
	"github.com/insanjati/fitbyte/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]model.UserResponse, error) {
	query := `SELECT name, email, preference, weight_unit, height_unit, weight, height, image_uri FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserResponse
	for rows.Next() {
		user := model.UserResponse{}
		err := rows.Scan(
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
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
