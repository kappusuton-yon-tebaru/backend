package build

import (
	"context"
	"encoding/json"

	"github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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

func (h *Handler) BuildImageHandler(msg amqp091.Delivery) {
	var body BuildRequestDTO

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		h.logger.Error("error occured while creating buidler pod", zap.Error(err))
		return
	}

	config := kubernetes.BuildImageDTO{
		Id:           body.Id,
		Dockerfile:   body.Dockerfile,
		Url:          body.Url,
		Destinations: body.Destinations,
	}

	ctx := context.Background()
	werr := h.service.BuildImage(ctx, config)
	if werr != nil {
		h.logger.Error("error occured while building image", zap.Error(err))
		return
	}
}
