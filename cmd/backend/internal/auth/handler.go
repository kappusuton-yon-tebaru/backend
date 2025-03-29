package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kappusuton-yon-tebaru/backend/internal/auth"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/httputils"
	"github.com/kappusuton-yon-tebaru/backend/internal/user"
	"github.com/kappusuton-yon-tebaru/backend/internal/validator"
)

type Handler struct {
	config    *config.Config
	service   *auth.Service
	validator *validator.Validator
}

func NewHandler(config *config.Config, service *auth.Service, validator *validator.Validator) *Handler {
	return &Handler{
		config,
		service,
		validator,
	}
}

// Register
//
//	@Router			/auth/register [post]
//	@Summary		Register
//	@Description	Register
//	@Tags			Auth
//	@Produce		json
//	@Param			request	body		RegisterRequest	true	"register request"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	httputils.ErrResponse
//	@Failure		500		{object}	httputils.ErrResponse
func (h *Handler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid body",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	credential := user.UserCredentialDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	werr := h.service.Register(ctx, credential)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.JSON(http.StatusCreated, AuthResponse{
		Message: "user created",
	})
}

// Login
//
//	@Router			/auth/login [post]
//	@Summary		Login
//	@Description	Login
//	@Tags			Auth
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"register request"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	httputils.ErrResponse
//	@Failure		500		{object}	httputils.ErrResponse
func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: "invalid body",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, httputils.ErrResponse{
			Message: strings.Join(h.validator.Translate(err), ", "),
		})
		return
	}

	credential := user.UserCredentialDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	token, werr := h.service.Login(ctx, credential)
	if werr != nil {
		ctx.JSON(httputils.ErrorResponseFromWErr(werr))
		return
	}

	ctx.SetCookie("session_token", token, h.config.SessionExpiresInSecond, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, AuthResponse{
		Message: "logged in",
	})
}

// Logout
//
//	@Router			/auth/logout [post]
//	@Summary		Logout
//	@Description	Logout
//	@Tags			Auth
//	@Success		200
func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("session_token", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, AuthResponse{
		Message: "logged out",
	})
}
