package kubernetes

import (
	"context"
	"io"

	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Pod struct {
	client v1.PodInterface
}

func (kube *Kubernetes) NewPodClient(namespace string) Pod {
	return Pod{
		client: kube.clientset.CoreV1().Pods(namespace),
	}
}

func (p Pod) Create(ctx context.Context, pod *apicorev1.Pod) (*apicorev1.Pod, error) {
	createdPod, err := p.client.Create(ctx, pod, apimetav1.CreateOptions{
		FieldManager: "system",
	})
	if err != nil {
		return nil, err
	}

	return createdPod, nil
}

func (p Pod) Delete(ctx context.Context, name string) error {
	err := p.client.Delete(ctx, name, apimetav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p Pod) Get(ctx context.Context, name string) (*apicorev1.Pod, error) {
	pod, err := p.client.Get(ctx, name, apimetav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (p Pod) GetLogString(ctx context.Context, name string) (string, error) {
	logs := p.client.GetLogs(name, &apicorev1.PodLogOptions{})
	log, err := logs.Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(log), err
}

func (p Pod) GetLogStream(ctx context.Context, name string) (io.ReadCloser, error) {
	logs := p.client.GetLogs(name, &apicorev1.PodLogOptions{})
	reader, err := logs.Stream(ctx)
	if err != nil {
		return nil, err
	}

	return reader, err
}
