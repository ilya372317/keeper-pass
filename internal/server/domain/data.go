package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
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

// KindsCanBeSimpleDeleted includes only data kinds which can be simple deleted from data storage.
// Not required any additional manipulation on deleting.
var KindsCanBeSimpleDeleted = []Kind{KindLoginPass, KindCreditCard, KindText, KindBinary}

// CryptoKeyLength length of general key.
const CryptoKeyLength = 32

// DataRepresentation interface for types which able to represent himself as string.
type DataRepresentation interface {
	GetInfo() string
}

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

// ToDataRepresentation parse union data type to specific by kind.
func (d Data) ToDataRepresentation() (DataRepresentation, error) {
	var result DataRepresentation
	switch d.Kind {
	case KindLoginPass:
		lp := LoginPass{}
		if err := json.Unmarshal([]byte(d.Metadata), &lp.Metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata for unmarshal to login pass data type: %w", err)
		}
		result = lp
	case KindCreditCard:
		cc := CreditCard{}
		if err := json.Unmarshal([]byte(d.Metadata), &cc.Metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata for unmarshal to credit card data type: %w", err)
		}
		result = cc
	case KindText:
		t := Text{}
		if err := json.Unmarshal([]byte(d.Metadata), &t.Metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata for unmarshal to text data type: %w", err)
		}
		result = t
	case KindBinary:
		b := Binary{}
		if err := json.Unmarshal([]byte(d.Metadata), &b.Metadata); err != nil {
			return nil, fmt.Errorf("invalid metadata for unmarshal to binary data type: %w", err)
		}
		result = b
	}
	if result == nil {
		return nil, fmt.Errorf("failed convert data to specific representation. Unknown kind")
	}

	return result, nil
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

// GetInfo parse LoginPass to string.
func (lp LoginPass) GetInfo() string {
	return lp.Metadata.URL
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

// GetInfo parse CreditCard to string.
func (cc CreditCard) GetInfo() string {
	return cc.Metadata.BankName
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

// GetInfo parse Text to string.
func (t Text) GetInfo() string {
	return t.Metadata.Info
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

// GetInfo parse Binary to string.
func (b Binary) GetInfo() string {
	return b.Metadata.Info
}

// GeneralData represent universal data representation for all types. Not Provide actual payload of data because
// it required for decrypt and may be slow and unnecessary for list information list.
type GeneralData struct {
	Info string `json:"info"`
	ID   int64  `json:"id"`
	Kind int8   `json:"kind"`
}
