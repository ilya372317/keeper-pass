package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid user password given")
)

// User entity which represent users.
type User struct {
	Email          string // user email.
	HashedPassword string // user hashed password.
	Salt           string // random generated salt for hash check.
}

func (u *User) IsPasswordCorrect(passwordHash string) bool {
	userPassWithSaltBytes := sha256.Sum256([]byte(passwordHash + u.Salt))
	userPassWithSalt := hex.EncodeToString(userPassWithSaltBytes[:])

	return u.HashedPassword == userPassWithSalt
}

func (u *User) SetHashedPassword(hashedPassword string) {
	passwordWithSaltBytes := sha256.Sum256([]byte(hashedPassword + u.Salt))
	passwordWithSalt := hex.EncodeToString(passwordWithSaltBytes[:])

	u.HashedPassword = passwordWithSalt
}
