package dto

type CreditCardMetadata struct {
	BankName string `json:"bank_name,omitempty"`
}

type SaveCreditCardDTO struct {
	Metadata   CreditCardMetadata `json:"-"`
	CardNumber string             `validate:"credit_card" json:"card_number"`
	Expiration string             `validate:"cardexp" json:"expiration"`
	CVV        int                `validate:"min=100,max=9999" json:"cvv"`
}
