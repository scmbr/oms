package dto

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
