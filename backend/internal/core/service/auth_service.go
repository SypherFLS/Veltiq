package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"veltiq/internal/core/domain"
	"veltiq/internal/core/ports"
)

type AuthService struct {
	users ports.UserRepository
	hasher ports.PasswordHasher
}

func NewAuthService(users ports.UserRepository, hasher ports.PasswordHasher) *AuthService {
	return &AuthService{
		users: users,
		hasher: hasher,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) error {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID: uuid.NewString(),
		Email: email,
		PasswordHash: hash,
		CreatedAt: time.Now(),
	}

	err = s.users.Create(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailTaken
		}
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	err = s.hasher.Compare(user.PasswordHash, password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return user.ID, nil
}
