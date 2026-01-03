package security

import (
	"github.com/v2code/b16/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BCryptPasswordHasher struct {
	cost int
}

func NewBCryptPasswordHasher(cost int) domain.PasswordHasher {
	return &BCryptPasswordHasher{
		cost: cost,
	}
}

func (h *BCryptPasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *BCryptPasswordHasher) Compare(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
