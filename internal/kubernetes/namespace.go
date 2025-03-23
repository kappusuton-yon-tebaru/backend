package kubernetes

import (
	"context"
	"fmt"

	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Namespace struct {
	client v1.NamespaceInterface
}

func (kube *Kubernetes) NewNamespaceClient() Namespace {
	return Namespace{
		client: kube.clientset.CoreV1().Namespaces(),
	}
}

func (ns Namespace) ListNamespaceByProjectId(ctx context.Context, projectId string) (*apicorev1.NamespaceList, error) {
	namespaces, err := ns.client.List(ctx, apimetav1.ListOptions{
		LabelSelector: fmt.Sprintf("project_id=%s", projectId),
	})
	if err != nil {
		return nil, err
	}

	return namespaces, nil
}

func (ns Namespace) GetNamespace(ctx context.Context, name string) (*apicorev1.Namespace, error) {
	namespace, err := ns.client.Get(ctx, name, apimetav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (ns Namespace) Create(ctx context.Context, projectId string, name string) (*apicorev1.Namespace, error) {
	dto := &apicorev1.Namespace{
		ObjectMeta: apimetav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"project_id": projectId,
			},
		},
	}

	namespace, err := ns.client.Create(ctx, dto, apimetav1.CreateOptions{
		FieldManager: SystemServiceAccount,
	})
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (ns Namespace) Delete(ctx context.Context, name string) error {
	err := ns.client.Delete(ctx, name, apimetav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
