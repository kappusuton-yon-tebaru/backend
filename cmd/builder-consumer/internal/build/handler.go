package build

import (
	"context"
	"encoding/json"

	"github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	logger  *logger.Logger
	service *build.Service
}

func NewHandler(logger *logger.Logger, service *build.Service) *Handler {
	return &Handler{
		logger,
		service,
	}
}

func (h *Handler) BuildImageHandler(msg amqp091.Delivery) *werror.WError {
	var body BuildRequestDTO

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		return werror.NewFromError(err)
	}

	config := kubernetes.BuildImageDTO{
		Dockerfile:   body.Dockerfile,
		Url:          body.Url,
		Destinations: body.Destinations,
		AppName:      body.AppName,
	}

	ctx := context.Background()
	werr := h.service.BuildImage(ctx, config)
	if werr != nil {
		return werr
	}

	return nil
}
