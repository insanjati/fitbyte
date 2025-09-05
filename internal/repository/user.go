package repository

import (
	"context"
	"fmt"
	"time"

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
	query := `SELECT name, email, preference, weightUnit, heightUnit, weight, height, imageUri FROM users WHERE id = $1`

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

func (r *UserRepository) UpdateUser(id uuid.UUID, user *model.UpdateUserRequest) (*model.UserResponse, error) {
	query := `UPDATE users 
	          SET name = $1, preference = $2, weightUnit = $3, heightUnit = $4, 
	              weight = $5, height = $6, imageUri = $7 
	          WHERE id = $8
	          RETURNING name, email, preference, weightUnit, heightUnit, weight, height, imageUri`

	var updated model.UserResponse
	err := r.db.QueryRow(query,
		user.Name,
		user.Preference,
		user.WeightUnit,
		user.HeightUnit,
		user.Weight,
		user.Height,
		user.ImageUri,
		id,
	).Scan(
		&updated.Name,
		&updated.Email,
		&updated.Preference,
		&updated.WeightUnit,
		&updated.HeightUnit,
		&updated.Weight,
		&updated.Height,
		&updated.ImageUri,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

// Repository for doing CRUD
func (r *UserRepository) RegisterNewUser(c context.Context, payload model.User) (model.User, error) {
	newId := uuid.New().String() // generate ID with UUID

	var user model.User //Assign user variable with model.User type

	//Query for insert new user data
	query := `INSERT INTO Users(id, email, password, created_at, updated_at) 
				VALUES($1, $2, $3, $4, $5) 
				RETURNING id, email`

	//Execute query that is expected to return at most one row and copy few results into variables
	err := r.db.QueryRowContext(c, query, newId, payload.Email, payload.Password, time.Now(), time.Now()).Scan(&user.ID, &user.Email)
	if err != nil { //error exists

		if c.Err() != nil { //if error from context exists

			return model.User{}, fmt.Errorf("context error: %w", c.Err()) //return context error only

		}
		return model.User{}, fmt.Errorf("operation failed: %w", err) // return operation failed error
	}

	return user, nil //return user and token
}

func (r *UserRepository) GetUserByEmail(c context.Context, email string) (model.User, error) {
	var user model.User

	query := `SELECT id, name, email, password FROM users WHERE email=$1`
	fmt.Println(email)
	err := r.db.QueryRowContext(c, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if c.Err() != nil {
			return model.User{}, c.Err()
		}
		return model.User{}, err
	}
	return user, nil
}
