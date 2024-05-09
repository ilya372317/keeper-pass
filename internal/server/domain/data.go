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

// LoginPassMetadata represent metadata of login pass data type.
type LoginPassMetadata struct {
	URL string
}

// LoginPass represent login pass data type.
type LoginPass struct {
	Metadata LoginPassMetadata
	Login    string
	Password string
	ID       int
}

// CreditCardMetadata represent metadata for credit card data type.
type CreditCardMetadata struct {
	BankName string `json:"bank_name,omitempty"`
}

// CreditCard represent credit card information data type.
type CreditCard struct {
	Metadata   CreditCardMetadata `json:"metadata,omitempty"`
	CardNumber string             `json:"card_number"`
	Expiration string             `json:"expiration"`
	CVV        int                `json:"cvv"`
	ID         int                `json:"id,omitempty"`
}

// TextMetadata represent metadata of text type.
type TextMetadata struct {
	Info string `json:"info"`
}

// Text represent text data type.
type Text struct {
	Metadata TextMetadata `json:"metadata,omitempty"`
	Data     string       `json:"data"`
	ID       int64        `json:"id"`
}

// BinaryMetadata represent metadata of binary type.
type BinaryMetadata struct {
	Info string `json:"info"`
}

// Binary represent binary data type.
type Binary struct {
	Metadata BinaryMetadata `json:"metadata,omitempty"`
	Data     []byte         `json:"data"`
	ID       int64          `json:"id"`
}

// GeneralData represent universal data representation for all types. Not Provide actual payload of data because
// it required for decrypt and may be slow and unnecessary for list information list.
type GeneralData struct {
	Info string `json:"info"`
	ID   int64  `json:"id"`
	Kind int8   `json:"kind"`
}
