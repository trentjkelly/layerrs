package authRepository

import (
	"golang.org/x/crypto/bcrypt"
	"context"
)

type PasswordRepository struct {}

func NewPasswordRepository() *PasswordRepository {
	passwordRepository := new(PasswordRepository)
	return passwordRepository
}

func (p *PasswordRepository) HashPassword(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p *PasswordRepository) CheckPassword(ctx context.Context, password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}