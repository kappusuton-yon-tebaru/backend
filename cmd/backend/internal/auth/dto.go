package auth

type RegisterReq struct {
	Email           string `json:"email"         validate:"required,email"`
	Password        string `json:"password"         validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
