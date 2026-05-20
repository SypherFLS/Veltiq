package postgres

import (
	"context"
	"errors"

	"veltiq/internal/core/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	record := userFromDomain(user)
	err := r.db.WithContext(ctx).Create(&record).Error
	if err != nil && isDuplicateKey(err) {
		return gorm.ErrDuplicatedKey
	}
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var record userRecord
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return userToDomain(record), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var record userRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return userToDomain(record), nil
}
