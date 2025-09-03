package model

import "github.com/google/uuid"

type AuthResponse struct{
	Email     string    `json:"email" db:"email"`
	Token     string    `json:"token"`
}

type AuthRequest struct{
	ID        uuid.UUID    `json:"id" db:"id"`	
	Email     string    `json:"email" db:"email"`
	Password  string	`json:"password" db:"password"`
}