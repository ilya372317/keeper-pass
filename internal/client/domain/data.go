package domain

import "errors"

var ErrUnauthenticated = errors.New("unauthenticated")

// Kind represent a data type.
type Kind int

const (
	KindLoginPass  Kind = iota // KindLoginPass represent a login with password data type.
	KindCreditCard             // KindCreditCard represent a credit card data type.
	KindText                   // KindText represent a text data type.
	KindBinary                 // KindBinary represent a binary data type.
)

var kindToString = map[Kind]string{
	KindLoginPass:  "login-pass",
	KindCreditCard: "credit-card",
	KindText:       "text",
	KindBinary:     "binary",
}

type ShortData struct {
	Info string `json:"info"`
	ID   int64  `json:"id"`
	Kind Kind   `json:"kind"`
}

func (sd ShortData) StringKind() string {
	return kindToString[sd.Kind]
}
