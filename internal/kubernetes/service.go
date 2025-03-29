package kubernetes

import (
	"context"

	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	accorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Service struct {
	client v1.ServiceInterface
}

func (kube *Kubernetes) NewServiceClient(namespace string) Service {
	return Service{
		client: kube.clientset.CoreV1().Services(namespace),
	}
}

func (s Service) Apply(ctx context.Context, svc *accorev1.ServiceApplyConfiguration) (*apicorev1.Service, error) {
	appliedService, err := s.client.Apply(ctx, svc, apimetav1.ApplyOptions{
		FieldManager: SystemServiceAccount,
	})
	if err != nil {
		return nil, err
	}

	return appliedService, nil
}

func (s Service) Get(ctx context.Context, name string) (*apicorev1.Service, error) {
	svc, err := s.client.Get(ctx, name, apimetav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s Service) Delete(ctx context.Context, name string) error {
	err := s.client.Delete(ctx, name, apimetav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
