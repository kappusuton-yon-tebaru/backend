package monitoring

import (
	"context"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/hub"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/streamer"
	apicorev1 "k8s.io/api/core/v1"
)

type Service struct {
	namespace string
	kube      *kubernetes.Kubernetes
	hub       *hub.Hub
}

func NewService(cfg *config.Config, kube *kubernetes.Kubernetes, hub *hub.Hub) *Service {
	return &Service{
		cfg.KubeNamespace,
		kube,
		hub,
	}
}

func (s *Service) GetPodLogs(ctx context.Context, name string) (*streamer.Streamer, error) {
	podClient := s.kube.NewPodClient(s.namespace)
	_, err := podClient.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	stream := s.hub.GetOrRegisterStreamer(name)

	go func() {
		ctx := context.Background()
		for {
			pod, err := podClient.Get(ctx, name)
			if err != nil {
				panic(err)
			}

			if pod.Status.Phase == apicorev1.PodRunning {
				break
			} else if pod.Status.Phase == apicorev1.PodFailed {
				return
			}
		}

		reader, err := podClient.GetLogStream(ctx, name)
		if err != nil {
			// log here
			return
		}
		defer reader.Close()

		stream.StreamFromReader(reader)
	}()

	return stream, nil
}
