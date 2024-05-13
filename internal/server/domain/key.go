package domain

import "time"

// Keys represent key for crypt and decrypt data keys.
type Keys struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Key       string    `db:"key"`
	Nonce     string    `db:"nonce"`
	ID        int       `db:"id"`
	IsCurrent bool      `db:"is_current"`
}
