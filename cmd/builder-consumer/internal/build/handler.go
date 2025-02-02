package build

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
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

func (h *Handler) BuildImageHandler(msg amqp091.Delivery) error {
	var body BuildRequestDTO

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		return fmt.Errorf("error while unmarshaling: %v", err)
	}

	config := kubernetes.BuildImageDTO{
		Dockerfile:   body.Dockerfile,
		Url:          body.Url,
		Destinations: body.Destinations,
		AppName:      body.AppName,
	}

	ctx := context.Background()
	err = h.service.BuildImage(ctx, config)
	if err != nil {
		return fmt.Errorf("error while create pod: %v", err)
	}

	return nil
}
