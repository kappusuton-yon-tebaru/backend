package deploy

import (
	"context"
	"encoding/json"

	sharedDeploy "github.com/kappusuton-yon-tebaru/backend/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *logger.Logger
	service *Service
}

func NewHandler(logger *logger.Logger, service *Service) *Handler {
	return &Handler{
		logger,
		service,
	}
}

func (h *Handler) DeployHandler(msg amqp091.Delivery) {
	var body sharedDeploy.DeployContext

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		h.logger.Error("error occured while parsing deploy context", zap.Error(err))
		return
	}

	h.logger.Info("consuming job", zap.String("job_id", body.Id))

	dto := kubernetes.DeployDTO{
		Id:           body.Id,
		ServiceName:  body.ServiceName,
		ImageUri:     body.ImageUri,
		Port:         body.Port,
		Namespace:    body.Namespace,
		Environments: body.Environments,
	}

	ctx := context.Background()
	werr := h.service.Deploy(ctx, dto)
	if werr != nil {
		h.logger.Error("error occured while deploying service", zap.Error(err))
		return
	}
}
