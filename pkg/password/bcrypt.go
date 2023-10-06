package password

import (
	"errors"
	"fmt"

	bcryptLib "golang.org/x/crypto/bcrypt"

	"github.com/taraslis453/solid-software-test/pkg/logging"
)

type bcrypt struct{}

// Check if implements the interface.
var _ Hasher = (*bcrypt)(nil)

func NewBcrypt(l logging.Logger) *bcrypt {
	return &bcrypt{}
}

func (bc *bcrypt) GenerateHashFromPassword(password string) (string, error) {
	hashedPassword, err := bcryptLib.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (bc *bcrypt) CompareHashAndPassword(opts *CompareHashAndPasswordOptions) (bool, error) {
	err := bcryptLib.CompareHashAndPassword([]byte(opts.Hashed), []byte(opts.Password))
	if errors.Is(err, bcryptLib.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check password correctness: %w", err)
	}

	return true, nil
}
