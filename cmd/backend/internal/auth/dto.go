package auth

type RegisterReq struct {
	Email           string `json:"username"         validate:"required,email"`
	Password        string `json:"password"         validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
