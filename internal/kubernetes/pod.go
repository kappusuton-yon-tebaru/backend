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
		FieldManager: SystemServiceAccount,
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

type PodLogBuilder struct {
	pod    string
	opt    *apicorev1.PodLogOptions
	client v1.PodInterface
}

func (p Pod) GetLog(pod string, opts ...PodLogOptions) PodLogBuilder {
	logOpt := &apicorev1.PodLogOptions{
		Timestamps: true,
	}

	for _, opt := range opts {
		opt(logOpt)
	}

	return PodLogBuilder{
		pod,
		logOpt,
		p.client,
	}
}

func (p PodLogBuilder) String(ctx context.Context) (string, error) {
	req := p.client.GetLogs(p.pod, p.opt)
	log, err := req.DoRaw(ctx)
	if err != nil {
		return "", err
	}

	return string(log), err
}

func (p PodLogBuilder) Stream(ctx context.Context) (io.ReadCloser, error) {
	p.opt.Follow = true
	req := p.client.GetLogs(p.pod, p.opt)
	reader, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}

	return reader, err
}
