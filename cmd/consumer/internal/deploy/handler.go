package deploy

import (
	"encoding/json"
	"fmt"

	sharedDeploy "github.com/kappusuton-yon-tebaru/backend/internal/deploy"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Handler struct {
	logger *logger.Logger
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		logger,
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

	fmt.Println(body)
}
