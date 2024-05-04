package dto

type CreditCardMetadata struct {
	BankName string
}

type SaveCreditCardDTO struct {
	Metadata   CreditCardMetadata
	CardNumber string `validate:"credit_card"`
	Expiration string `validate:"cardexp"`
	CVV        int    `validate:"min=100,max=9999"`
}
