package postgres

import (
	"time"

	"veltiq/internal/core/domain"
)

type userRecord struct {
	ID string `gorm:"column:id;primaryKey;type:uuid"`
	Email string `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string `gorm:"column:password_hash;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (userRecord) TableName() string {
	return "users"
}

func userFromDomain(u *domain.User) userRecord {
	return userRecord{
		ID: u.ID,
		Email: u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt: u.CreatedAt,
	}
}

func userToDomain(r userRecord) *domain.User {
	return &domain.User{
		ID: r.ID,
		Email: r.Email,
		PasswordHash: r.PasswordHash,
		CreatedAt: r.CreatedAt,
	}
}
