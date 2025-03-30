package podevent

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/cmd/podlogger/internal/podwatcher"
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
)

type Handler struct {
	logger   *logger.Logger
	kube     *kubernetes.Kubernetes
	watchers map[string]chan struct{}
	sync.Mutex
}

func NewHandler(kube *kubernetes.Kubernetes, logger *logger.Logger) *Handler {
	return &Handler{
		logger:   logger,
		kube:     kube,
		watchers: make(map[string]chan struct{}),
	}
}

func (h *Handler) PodCreated(pod *apicorev1.Pod) {
	container, ok := pod.Labels["watchlog.container"]
	if !ok {
		container = ""
	}

	attributes, ok := pod.Labels["watchlog.attributes"]
	if !ok {
		attributes = ""
	}

	attrs := map[string]string{}
	for _, key := range strings.Split(attributes, ".") {
		val, ok := pod.Labels[key]
		if !ok {
			continue
		}

		attrs[key] = val
	}

	term := make(chan struct{})

	go func() {
		podWatcher := podwatcher.NewPodWatcher(h.kube, pod.Namespace, pod.Name, container)

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
						h.logger.Error("error occured while getting pod log", zap.String("pod", pod.Name), zap.Error(err))
						time.Sleep(ReconnectTimeout)
						continue
					}

					time.Sleep(ReconnectTimeout)
					h.logger.Info("trying to access pod log", zap.String("pod", pod.Name))
				}
			}
		}()

		for chunk := range chunkLogCh {
			for _, log := range chunk {
				maps.Copy(log.Attribute, attrs)

				bs, _ := json.Marshal(log)
				fmt.Println(string(bs))
			}
		}
	}()

	h.logger.Info("Registered pod", zap.String("namespace", pod.Namespace), zap.String("pod", pod.Name), zap.String("container", container))

	h.Lock()
	h.watchers[pod.Name] = term
	h.Unlock()
}

func (h *Handler) PodUpdated(oldPod *apicorev1.Pod, newPod *apicorev1.Pod) {}

func (h *Handler) PodDeleted(pod *apicorev1.Pod) {
	h.logger.Info("Unregistered pod", zap.String("namespace", pod.Namespace), zap.String("pod", pod.Name))

	h.Lock()
	close(h.watchers[pod.Name])
	delete(h.watchers, pod.Name)
	h.Unlock()
}
