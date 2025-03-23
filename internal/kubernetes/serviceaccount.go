package kubernetes

import (
	"context"

	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ServiceAccount struct {
	client    v1.ServiceAccountInterface
	namespace string
}

func (kube *Kubernetes) NewServiceAccountClient(namespace string) ServiceAccount {
	return ServiceAccount{
		client:    kube.clientset.CoreV1().ServiceAccounts(namespace),
		namespace: namespace,
	}
}

func (s ServiceAccount) Create(ctx context.Context, name string) error {
	serviceAccount := &apicorev1.ServiceAccount{
		ObjectMeta: apimetav1.ObjectMeta{
			Name:      name,
			Namespace: s.namespace,
		},
	}

	_, err := s.client.Create(ctx, serviceAccount, apimetav1.CreateOptions{
		FieldManager: SystemServiceAccount,
	})
	if err != nil {
		return err
	}

	return nil
}
