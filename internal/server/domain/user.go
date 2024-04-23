package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const saltLength = 32

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid user password given")
)

type CtxUserKey struct{}

// User entity which represent users.
type User struct {
	CreatedAT      time.Time `db:"created_at"`      // created at date.
	UpdatedAT      time.Time `db:"updated_at"`      // updated at date.
	Email          string    `db:"email"`           // user email.
	HashedPassword string    `db:"hashed_password"` // user hashed password.
	Salt           string    `db:"salt"`            // random generated salt for hash check.
	ID             uint      `db:"id"`              // identifier
}

func (u *User) IsPasswordCorrect(password string) bool {
	userPassWithSaltBytes := sha256.Sum256([]byte(password + u.Salt))
	userPassWithSalt := hex.EncodeToString(userPassWithSaltBytes[:])

	return u.HashedPassword == userPassWithSalt
}

// SetHashedPassword from given user password generate hash and set in HashedPassword field.
func (u *User) SetHashedPassword(password string) {
	passwordWithSaltBytes := sha256.Sum256([]byte(password + u.Salt))
	passwordWithSalt := hex.EncodeToString(passwordWithSaltBytes[:])

	u.HashedPassword = passwordWithSalt
}

// GenerateSalt generate salt for make password hash more strong.
func (u *User) GenerateSalt() error {
	saltBytes := make([]byte, saltLength)
	if _, err := rand.Read(saltBytes); err != nil {
		return fmt.Errorf("failed generate salt for user: %w", err)
	}

	salt := hex.EncodeToString(saltBytes)
	u.Salt = salt

	return nil
}
