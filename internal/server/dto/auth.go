package dto

type RegisterDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=3,max=255"`
}

type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=3,max=255"`
}

type JWTClaimsDTO struct {
	Email string
}
