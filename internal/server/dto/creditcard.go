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

type UpdateCreditCardDTO struct {
	Metadata   *CreditCardMetadata `json:"_" validate:"omitnil"`
	CardNumber *string             `json:"card_number" validate:"omitnil,credit_card"`
	Expiration *string             `json:"expiration" validate:"omitnil,cardexp"`
	CVV        *int32              `json:"cvv" validate:"omitnil,min=100,max=9999"`
	ID         int
}
