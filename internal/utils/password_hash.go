package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	EncryptPassword(password string) (string, error)
	ComparePasswordHash(passwordHash string, plainPassword string) error
}

type passwordHasher struct{}

func NewPasswordHasher() PasswordHasher {
	return &passwordHasher{}
}

func (ph *passwordHasher) EncryptPassword(password string) (string, error) {
	passwordHash, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if error != nil {
		return "", error
	}
	return string(passwordHash), nil
}

func (ph *passwordHasher) ComparePasswordHash(passwordHash string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
}
