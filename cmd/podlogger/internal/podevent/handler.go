package podevent

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/cmd/podlogger/internal/podwatcher"
	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/logging"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"go.uber.org/zap"
	apicorev1 "k8s.io/api/core/v1"
)

const (
	BufferSize        = 64
	ReconnectTimeout  = 3 * time.Second
	DebouncedInterval = 3 * time.Second

	InsertRetryTimeout = 1 * time.Second
	MaximumInsertRetry = 5
)

type Handler struct {
	mode           enum.PodLoggerMode
	logger         *logger.Logger
	loggingService *logging.Service
	kube           *kubernetes.Kubernetes
	watchers       map[string]chan struct{}
	sync.Mutex
}

func NewHandler(config *config.Config, kube *kubernetes.Kubernetes, logger *logger.Logger, loggingService *logging.Service) *Handler {
	return &Handler{
		mode:           config.PodLogger.Mode,
		logger:         logger,
		loggingService: loggingService,
		kube:           kube,
		watchers:       make(map[string]chan struct{}),
	}
}

func (h *Handler) PodCreated(pod *apicorev1.Pod) {
	container, ok := pod.Labels["watchlog.container"]
	if !ok {
		container = ""
	}

	attributes := h.podAttributes(pod)

	term := make(chan struct{})

	go func() {
		podWatcher := podwatcher.NewPodWatcher(h.kube, h.logger, pod, container)

		ch := make(chan logging.InsertLogDTO, BufferSize)
		chunkLogCh := utils.DebouncerChannel(ch, DebouncedInterval, BufferSize)

		go func() {
			for {
				select {
				case <-term:
					close(ch)
					return
				default:
					if err := podWatcher.WatchLog(context.Background(), ch); err != nil {
						h.logger.Error("error occured while getting pod log", zap.String("pod", pod.Name), zap.Duration("next_retry", ReconnectTimeout), zap.Error(err))
						time.Sleep(ReconnectTimeout)
						continue
					}

					time.Sleep(ReconnectTimeout)
				}
			}
		}()

		for chunk := range chunkLogCh {
			for _, log := range chunk {
				maps.Copy(log.Attribute, attributes)
			}

			h.insertLog(pod, chunk)
		}
	}()

	h.Lock()
	h.watchers[pod.Name] = term
	h.Unlock()

	h.logger.Info("Registered pod", zap.String("namespace", pod.Namespace), zap.String("pod", pod.Name), zap.String("container", container))
}

func (h *Handler) PodUpdated(oldPod *apicorev1.Pod, newPod *apicorev1.Pod) {
	restarts := map[string]int32{}
	for _, status := range newPod.Status.ContainerStatuses {
		restarts[status.Name] = status.RestartCount
	}

	for _, status := range oldPod.Status.ContainerStatuses {
		restarts[status.Name] -= status.RestartCount
	}

	logs := []logging.InsertLogDTO{}
	for name, value := range restarts {
		if value > 0 {
			terminated := (*apicorev1.ContainerStateTerminated)(nil)
			for _, condition := range newPod.Status.ContainerStatuses {
				if condition.Name == name {
					terminated = condition.State.Terminated
				}
			}

			logs = append(logs, logging.InsertLogDTO{
				Log:       fmt.Sprintf("[system] %s service terminated with exit code: %d, reason: %s", name, terminated.ExitCode, terminated.Reason),
				Attribute: h.podAttributes(newPod),
				Timestamp: time.Now(),
			})
		}
	}

	if len(logs) > 0 {
		h.insertLog(newPod, logs)
	}
}

func (h *Handler) PodDeleted(pod *apicorev1.Pod) {
	h.Lock()
	close(h.watchers[pod.Name])
	delete(h.watchers, pod.Name)
	h.Unlock()

	h.logger.Info("Unregistered pod", zap.String("namespace", pod.Namespace), zap.String("pod", pod.Name))
}

func (h *Handler) podAttributes(pod *apicorev1.Pod) map[string]string {
	attributes, ok := pod.Labels["watchlog.attributes"]
	if !ok {
		attributes = ""
	}

	attrs := map[string]string{}
	for _, key := range strings.Split(attributes, ".") {
		val, ok := pod.Labels[key]
		if ok {
			attrs[key] = val
		}
	}

	return attrs
}

func (h *Handler) insertLog(pod *apicorev1.Pod, chunk []logging.InsertLogDTO) {
	switch h.mode {
	case enum.PodLoggerModeMongoDb:
		werr := h.loggingService.BatchInsertLog(context.Background(), chunk)
		if werr != nil {
			h.logger.Info(
				"error occured while inserting logs",
				zap.Error(werr.Err),
				zap.String("pod", pod.Name),
				zap.Duration("next_retry", InsertRetryTimeout),
			)

			go h.retryInsertion(pod.Name, chunk)
		}
	default:
		for _, log := range chunk {
			h.logger.Info(
				log.Log,
				zap.Time("timestamp", log.Timestamp),
				zap.Any("attribute", log.Attribute),
			)
		}
	}
}

func (h *Handler) retryInsertion(pod string, chunk []logging.InsertLogDTO) {
	timeout := InsertRetryTimeout

	for range MaximumInsertRetry {
		time.Sleep(timeout)

		werr := h.loggingService.BatchInsertLog(context.Background(), chunk)
		if werr != nil {
			h.logger.Info(
				"error occured while inserting logs",
				zap.Error(werr.Err),
				zap.String("pod", pod),
				zap.Duration("next_retry", InsertRetryTimeout),
			)
		}

		timeout *= 2
	}
}
