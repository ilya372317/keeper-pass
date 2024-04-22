package dto

type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=3,max=255"`
}
