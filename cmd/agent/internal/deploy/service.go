package deploy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type Service struct {
	kube            *kubernetes.Kubernetes
	resourceService *resource.Service
}

func NewService(kube *kubernetes.Kubernetes, resourceService *resource.Service) *Service {
	return &Service{
		kube,
		resourceService,
	}
}

func (s *Service) DeleteDeployment(ctx context.Context, dto DeleteDeploymentRequest) *werror.WError {
	project, werr := s.resourceService.GetResourceByID(ctx, dto.ProjectId)
	if werr != nil {
		return werr
	}

	name := deployenv.GetNamespaceName(project.ResourceName, dto.DeploymentEnv)

	deployClient := s.kube.NewDeploymentClient(name)
	err := deployClient.Delete(ctx, fmt.Sprintf("%s-deployment", dto.ServiceName))
	if apierrors.IsNotFound(err) {
		return werror.NewFromError(err).SetCode(http.StatusBadRequest).SetMessage("deployment not found")
	} else if err != nil {
		return werror.NewFromError(err)
	}

	serviceClient := s.kube.NewServiceClient(name)
	err = serviceClient.Delete(ctx, fmt.Sprintf("%s-service", dto.ServiceName))
	if apierrors.IsNotFound(err) {
		return werror.NewFromError(err).SetCode(http.StatusBadRequest).SetMessage("service not found")
	} else if err != nil {
		return werror.NewFromError(err)
	}

	return nil
}
