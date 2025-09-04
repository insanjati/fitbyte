package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       *string   `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Preference string    `json:"preference" db:"preference"`
	WeightUnit string    `json:"weightUnit" db:"weightUnit"`
	HeightUnit string    `json:"heightUnit" db:"heightUnit"`
	Weight     int       `json:"weight" db:"weight"`
	Height     int       `json:"height" db:"height"`
	ImageUri   string    `json:"imageUri" db:"imageUri"`
	Password   string    `json:"password" db:"password"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type UserResponse struct {
	Name       *string  `json:"name"`
	Email      string   `json:"email"`
	Preference *string  `json:"preference"`
	WeightUnit *string  `json:"weightUnit"`
	HeightUnit *string  `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
	ImageUri   *string  `json:"imageUri"`
}

type UpdateUserRequest struct {
	Name       *string  `json:"name"`
	Preference *string  `json:"preference" validate:"required,oneof=CARDIO WEIGHT"`
	WeightUnit *string  `json:"weightUnit" validate:"required,oneof=KG LBS"`
	HeightUnit *string  `json:"heightUnit" validate:"required,oneof=CM INCH"`
	Weight     *float64 `json:"weight" validate:"required,gte=10,lte=1000"`
	Height     *float64 `json:"height" validate:"required,gte=3,lte=250"`
	ImageUri   *string  `json:"imageUri"`
}
