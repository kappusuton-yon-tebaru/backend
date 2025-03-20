package deployenv

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/logger"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type Service struct {
	kube            *kubernetes.Kubernetes
	resourceService *resource.Service
	logger          *logger.Logger
}

func NewService(kube *kubernetes.Kubernetes, resourceService *resource.Service, logger *logger.Logger) *Service {
	return &Service{
		kube,
		resourceService,
		logger,
	}
}

func (s *Service) ListDeploymentEnv(ctx context.Context, projectId string) ([]string, *werror.WError) {
	project, werr := s.resourceService.GetResourceByID(ctx, projectId)
	if werr != nil {
		return nil, werr
	}

	nsClient := s.kube.NewNamespaceClient()
	listResult, err := nsClient.ListNamespaceByProjectId(ctx, projectId)
	if err != nil {
		s.logger.Error("error occured while listing namespace", zap.String("project_id", projectId), zap.Error(err))
		return nil, werror.NewFromError(err)
	}

	namespaces := []string{}
	for _, namespace := range listResult.Items {
		namespaces = append(namespaces, strings.TrimPrefix(namespace.Name, fmt.Sprintf("%s-", project.ResourceName)))
	}

	return namespaces, nil
}

func (s *Service) CreateDeploymentEnv(ctx context.Context, dto ModifyDeploymentEnvDTO) *werror.WError {
	project, werr := s.resourceService.GetResourceByID(ctx, dto.ProjectId)
	if werr != nil {
		return werr
	}

	name := GetNamespaceName(project.ResourceName, dto.Name)

	nsClient := s.kube.NewNamespaceClient()
	_, err := nsClient.Create(ctx, dto.ProjectId, name)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return werror.NewFromError(err).SetMessage("deployment environment already exist on this project").SetCode(http.StatusBadRequest)
	} else if err != nil {
		s.logger.Error("error occured while creating namespace", zap.String("namespace", name), zap.Error(err))
		return werror.NewFromError(err)
	}

	saClient := s.kube.NewServiceAccountClient(name)
	err = saClient.Create(ctx, "system")
	if err != nil {
		s.logger.Error("error occured while creating service account", zap.String("namespace", project.ResourceName), zap.Error(err))
		return werror.NewFromError(err)
	}

	return nil
}

func (s *Service) DeleteDeploymentEnv(ctx context.Context, dto ModifyDeploymentEnvDTO) *werror.WError {
	project, werr := s.resourceService.GetResourceByID(ctx, dto.ProjectId)
	if werr != nil {
		return werr
	}

	name := GetNamespaceName(project.ResourceName, dto.Name)

	nsClient := s.kube.NewNamespaceClient()
	err := nsClient.Delete(ctx, name)
	if err != nil {
		s.logger.Error("error occured while deleting namespace", zap.String("namespace", name), zap.Error(err))
		return werror.NewFromError(err)
	}

	for {
		_, err := nsClient.GetNamespace(ctx, name)
		if err != nil && !apierrors.IsNotFound(err) {
			s.logger.Error("error occured while waiting for namespace to be deleted", zap.String("namespace", name), zap.Error(err))
			return werror.NewFromError(err)
		} else if apierrors.IsNotFound(err) {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (s *Service) DeleteDeployment(ctx context.Context, dto DeleteDeploymentDTO) *werror.WError {
	return nil
}

func GetNamespaceName(projectName, envName string) string {
	return fmt.Sprintf("%s-%s", projectName, envName)
}
