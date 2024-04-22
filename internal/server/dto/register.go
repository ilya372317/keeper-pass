package dto

type RegisterDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=3,max=255"`
}
