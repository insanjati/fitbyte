package model

import "time"

type User struct {
	ID         int       `json:"id" db:"id"`
	Name       *string   `json:"name" db:"name"`
	Email      string    `json:"email" db:"email" binding:"required,email"`
	Password   string    `json:"password" db:"password" binding:"required,min=6"`
	Preference *string   `json:"preference" db:"preference"`
	WeightUnit *string   `json:"weightUnit" db:"weight_unit"`
	HeightUnit *string   `json:"heightUnit" db:"height_unit"`
	Weight     *int      `json:"weight" db:"weight"`
	Height     *int      `json:"height" db:"height"`
	ImageUri   *string   `json:"imageUri" db:"image_uri"`
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
