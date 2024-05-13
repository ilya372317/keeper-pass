package domain

import (
	"errors"
	"fmt"
)

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

var AliasToKind = reverseKindToString(kindToString)

func reverseKindToString(m map[Kind]string) map[string]Kind {
	res := make(map[string]Kind, len(m))
	for k, a := range m {
		res[a] = k
	}

	return res
}

type ShortData struct {
	Info string `json:"info"`
	ID   int64  `json:"id"`
	Kind Kind   `json:"kind"`
}

func (sd ShortData) StringKind() string {
	return kindToString[sd.Kind]
}

type ShowAble interface {
	ToString() string
}

type CreditCard struct {
	BankName   string
	CardNumber string
	Exp        string
	Code       int
	ID         int
}

func (cc CreditCard) ToString() string {
	return fmt.Sprintf("ID : %d\nBANK NAME: %s\nCARD NUMBER: %s\nEXP:%s\nCODE:%d\n",
		cc.ID, cc.BankName, cc.CardNumber, cc.Exp, cc.Code)
}

type LoginPass struct {
	URL      string
	Login    string
	Password string
	ID       int
}

func (lp LoginPass) ToString() string {
	return fmt.Sprintf("ID: %d\nURL: %s\nLOGIN: %s\nPASSWORD: %s\n",
		lp.ID, lp.URL, lp.Login, lp.Password)
}

type Text struct {
	Info string
	Data string
	ID   int
}

func (t Text) ToString() string {
	return fmt.Sprintf("ID: %d\nINFO: %s\nDATA: %s\n",
		t.ID, t.Info, t.Data)
}

type Binary struct {
	Info string
	Data []byte
	ID   int
}

func (b Binary) ToString() string {
	return fmt.Sprintf("ID: %d\nINFO: %s\nDATA: %s\n",
		b.ID, b.Info, b.Data)
}
