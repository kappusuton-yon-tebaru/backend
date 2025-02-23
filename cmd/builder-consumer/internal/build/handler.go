package build

import (
	"context"
	"encoding/json"

	sharedBuild "github.com/kappusuton-yon-tebaru/backend/internal/build"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/regproviders"
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

func (h *Handler) BuildImageHandler(msg amqp091.Delivery) {
	var body sharedBuild.BuildContext

	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		h.logger.Error("error occured while parsing build context", zap.Error(err))
		return
	}

	h.logger.Info("consuming job", zap.String("job_id", body.Id))

	credentialMap, ok := body.RegistryCredential.(map[string]any)
	if !ok {
		h.logger.Error("error occured while parsing build context", zap.Error(err))
		return
	}

	credential, werr := regproviders.ParseCredential(body.RegistryType, credentialMap)
	if werr != nil {
		h.logger.Error("error occured while parsing build context", zap.Error(err))
		return
	}

	config := kubernetes.BuildImageDTO{
		Id:           body.Id,
		Dockerfile:   body.Dockerfile,
		RepoUrl:      body.RepoUrl,
		RepoRoot:     body.RepoRoot,
		Destinations: []string{body.Destination},
		Credential:   credential,
	}

	ctx := context.Background()
	werr = h.service.BuildImage(ctx, config)
	if werr != nil {
		h.logger.Error("error occured while building image", zap.Error(err))
		return
	}
}
