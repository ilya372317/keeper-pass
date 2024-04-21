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
	Email          string `db:"email"`           // user email.
	HashedPassword string `db:"hashed_password"` // user hashed password.
	Salt           string `db:"salt"`            // random generated salt for hash check.
	ID             uint   `db:"id"`              // identifier
}

func (u *User) IsPasswordCorrect(password string) bool {
	userPassWithSaltBytes := sha256.Sum256([]byte(password + u.Salt))
	userPassWithSalt := hex.EncodeToString(userPassWithSaltBytes[:])

	return u.HashedPassword == userPassWithSalt
}

func (u *User) SetHashedPassword(hashedPassword string) {
	passwordWithSaltBytes := sha256.Sum256([]byte(hashedPassword + u.Salt))
	passwordWithSalt := hex.EncodeToString(passwordWithSaltBytes[:])

	u.HashedPassword = passwordWithSalt
}
