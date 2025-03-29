package auth

type RegisterRequest struct {
	Email           string `json:"email"         validate:"required,email"`
	Password        string `json:"password"         validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Message string `json:"message"`
}
