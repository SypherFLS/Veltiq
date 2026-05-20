package service

import (
	"context"
	"time"

	"veltiq/internal/core/ports"
	"veltiq/internal/core/domain"
	"github.com/google/uuid"
)

type AuthService struct {
	users   ports.UserRepository
	hasher  ports.PasswordHasher
	tokens  ports.TokenManager
}

func NewAuthService(users ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenManager) *AuthService {
	return &AuthService{
		users:  users,
		hasher: hasher,
		tokens: tokens,
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
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = s.hasher.Compare(user.PasswordHash, password)
	if err != nil {
		return "", err
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}