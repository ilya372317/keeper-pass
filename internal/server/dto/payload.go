package dto

type LoginPassPayload struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type FilePayload struct {
	FileName string `json:"file_name" validate:"file,required"`
	MimeType string `json:"mime_type" validate:"required"`
}

type CreditCardPayload struct {
	CardNumber string `json:"card_number" validate:"required,credit_card"`
	Expiration string `json:"expiration" validate:"required"`
	CVV        string `json:"cvv" validate:"required"`
}
