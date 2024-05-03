package domain

import (
	"fmt"
	"time"
)

var ErrPayloadNotValid = fmt.Errorf("payload not valid")

// Kind represent a data type.
type Kind int

const (
	KindLoginPass  Kind = iota // KindLoginPass represent a login with password data type.
	KindFile                   // KindFile represent a file data type.
	KindCreditCard             // KindCreditCard represent a credit card data type.
)

const CryptoKeyLength = 32

// Data represent a storing data.
type Data struct {
	CreatedAt          time.Time `db:"created_at"`       // CreatedAt creation date.
	UpdatedAt          time.Time `db:"updated_at"`       // UpdatedAt // update date.
	Payload            string    `db:"payload"`          // Payload of the data.
	Metadata           string    `db:"metadata"`         // Metadata of the data.
	PayloadNonce       string    `db:"payload_nonce"`    // PayloadNonce of the data.
	CryptoKeyNonce     string    `db:"crypto_key_nonce"` // CryptoKeyNonce of the data.
	CryptoKey          string    `db:"crypto_key"`       // CryptoKey of the data.
	ID                 int       `db:"id"`               // ID of the data.
	UserID             uint      `db:"user_id"`          // UserID of the data.
	Kind               Kind      `db:"kind"`             // Kind of the data.
	IsPayloadDecrypted bool      // IsPayloadDecrypted flag indicates is payload encrypted or decrypted
}
