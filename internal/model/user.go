package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `json:"id" db:"id"`	
	Name      *string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`

	Preference string   `json:"preference" db:"preference"`
	WeightUnit string   `json:"weightUnit" db:"weightUnit"`
	HeightUnit string   `json:"heightUnit" db:"heightUnit"`
	
	Weight 	   int      `json:"weight" db:"weight"` 
	Height     int 		`json:"height" db:"height"`
	
	ImageUri	string  `json:"imageUri" db:"imageUri"`

	Password  string	`json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
