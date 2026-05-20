package security

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordManager struct{}

func (p *PasswordManager) Hash(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (p *PasswordManager) Compare(hash string, password string,) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}