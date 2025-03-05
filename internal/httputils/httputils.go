package httputils

import (
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func ErrorResponseFromWErr(werr *werror.WError) (int, ErrResponse) {
	return werr.GetCodeOr(http.StatusInternalServerError), ErrResponse{
		Message: werr.GetMessageOr("internal server error"),
	}
}
