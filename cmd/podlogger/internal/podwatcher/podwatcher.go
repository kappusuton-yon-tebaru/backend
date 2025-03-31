package podwatcher

import (
	"bufio"
	"context"
	"strings"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logging"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodWatcher struct {
	kube             *kubernetes.Kubernetes
	lastLogTimestamp *apimetav1.Time
	namespace        string
	pod              string
	container        string
}

func NewPodWatcher(kube *kubernetes.Kubernetes, pod *apicorev1.Pod, container string) *PodWatcher {
	return &PodWatcher{
		kube,
		(*apimetav1.Time)(nil),
		pod.Namespace,
		pod.Name,
		container,
	}
}

func (w *PodWatcher) WatchLog(ctx context.Context, out chan<- logging.InsertLogDTO) error {
	podClient := w.kube.NewPodClient(w.namespace)
	reader, err := podClient.GetLog(w.pod, kubernetes.WithContainer(w.container)).Stream(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		logs := strings.SplitN(line, " ", 2)
		timestamp, _ := time.Parse(time.RFC3339, logs[0])

		if w.lastLogTimestamp == nil || timestamp.UnixNano() > w.lastLogTimestamp.UnixNano() {
			w.lastLogTimestamp = utils.Pointer(apimetav1.NewTime(timestamp))
			msg := ""
			if len(logs) >= 2 {
				msg = logs[1]
			}

			out <- logging.InsertLogDTO{
				Timestamp: timestamp,
				Log:       msg,
				Attribute: map[string]string{},
			}
		}
	}

	return nil
}
