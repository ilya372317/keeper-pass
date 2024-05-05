package domain

import (
	"fmt"
	"time"
)

var (
	ErrPayloadNotValid       = fmt.Errorf("payload not valid")
	ErrNotSupportedOperation = fmt.Errorf("this method can`t get data of this kind")
)

// Kind represent a data type.
type Kind int

const (
	KindLoginPass  Kind = iota // KindLoginPass represent a login with password data type.
	KindCreditCard             // KindCreditCard represent a credit card data type.
	KindText                   // KindText represent a text data type.
	KindBinary                 // KindBinary represent a binary data type.
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

type LoginPassMetadata struct {
	URL string
}

type LoginPassData struct {
	Metadata LoginPassMetadata
	Login    string
	Password string
	ID       int
}

type CreditCardMetadata struct {
	BankName string `json:"bank_name,omitempty"`
}

type CreditCardData struct {
	Metadata   CreditCardMetadata `json:"metadata,omitempty"`
	CardNumber string             `json:"card_number"`
	Expiration string             `json:"expiration"`
	CVV        int                `json:"cvv"`
	ID         int                `json:"id,omitempty"`
}

type TextMetadata struct {
	Info string `json:"info"`
}

type Text struct {
	Metadata TextMetadata `json:"metadata,omitempty"`
	Info     string       `json:"info"`
	ID       int64        `json:"id"`
}

type BinaryMetadata struct {
	Info string `json:"info"`
}

type Binary struct {
	Metadata BinaryMetadata `json:"metadata,omitempty"`
	Info     []byte         `json:"info"`
	ID       int64          `json:"id"`
}
