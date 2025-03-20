package kubernetes

import (
	"context"

	apiappsv1 "k8s.io/api/apps/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	acappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type Deployment struct {
	client v1.DeploymentInterface
}

func (kube *Kubernetes) NewDeploymentClient(namespace string) Deployment {
	return Deployment{
		client: kube.clientset.AppsV1().Deployments(namespace),
	}
}

func (d Deployment) Apply(ctx context.Context, deployment *acappsv1.DeploymentApplyConfiguration) (*apiappsv1.Deployment, error) {
	appliedDeployment, err := d.client.Apply(ctx, deployment, apimetav1.ApplyOptions{
		FieldManager: "system",
	})
	if err != nil {
		return nil, err
	}

	return appliedDeployment, nil
}

func (d Deployment) Get(ctx context.Context, name string) (*apiappsv1.Deployment, error) {
	deployment, err := d.client.Get(ctx, name, apimetav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (d Deployment) GetCondition(deployment *apiappsv1.Deployment) DeploymentCondition {
	deployCondition := DeploymentCondition{}

	for _, condition := range deployment.Status.Conditions {
		switch condition.Type {
		case apiappsv1.DeploymentAvailable:
			deployCondition.Available = &condition
		case apiappsv1.DeploymentProgressing:
			deployCondition.Progressing = &condition
		case apiappsv1.DeploymentReplicaFailure:
			deployCondition.ReplicaFailure = &condition
		}
	}

	return deployCondition
}

func (d Deployment) Delete(ctx context.Context, name string) error {
	err := d.client.Delete(ctx, name, apimetav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
