package monitoring

import (
	"context"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/hub"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/streamer"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type Service struct {
	namespace string
	kube      *kubernetes.Kubernetes
	hub       *hub.Hub
	logger    *logger.Logger
}

func NewService(cfg *config.Config, kube *kubernetes.Kubernetes, hub *hub.Hub, logger *logger.Logger) *Service {
	return &Service{
		cfg.KubeNamespace,
		kube,
		hub,
		logger,
	}
}

func (s *Service) GetPodLogs(ctx context.Context, name string) (*streamer.Streamer, *werror.WError) {
	podClient := s.kube.NewPodClient(s.namespace)
	_, err := podClient.Get(ctx, name)
	if errors.IsNotFound(err) {
		return nil, werror.NewFromError(err).SetCode(http.StatusNotFound).SetMessage("pod not found")
	} else if err != nil {
		s.logger.Error("error occured while getting pod", zap.String("name", name), zap.Error(err))
		return nil, werror.NewFromError(err)
	}

	stream := s.hub.GetOrRegisterStreamer(name)

	go func() {
		ctx := context.Background()
		for {
			pod, err := podClient.Get(ctx, name)
			if err != nil {
				s.logger.Error("error occured while getting pod", zap.String("name", name), zap.Error(err))
				return
			}

			if pod.Status.Phase == apicorev1.PodRunning {
				break
			} else if pod.Status.Phase == apicorev1.PodFailed {
				return
			}
		}

		reader, err := podClient.GetLogStream(ctx, name, "kaniko")
		if err != nil {
			s.logger.Error("error occured while getting pod log reader", zap.String("name", name), zap.Error(err))
			return
		}
		defer func() {
			if err := reader.Close(); err != nil {
				s.logger.Error("error occured while closing log reader", zap.Error(err))
			}
		}()

		stream.StreamFromReader(reader)
	}()

	return stream, nil
}
