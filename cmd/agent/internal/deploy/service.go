package deploy

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/kappusuton-yon-tebaru/backend/internal/deployenv"
	"github.com/kappusuton-yon-tebaru/backend/internal/enum"
	"github.com/kappusuton-yon-tebaru/backend/internal/kubernetes"
	"github.com/kappusuton-yon-tebaru/backend/internal/models"
	"github.com/kappusuton-yon-tebaru/backend/internal/query"
	"github.com/kappusuton-yon-tebaru/backend/internal/resource"
	"github.com/kappusuton-yon-tebaru/backend/internal/utils"
	"github.com/kappusuton-yon-tebaru/backend/internal/werror"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type PaginatedDeployment models.Paginated[models.Deployment]

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

func (s *Service) ListDeployment(ctx context.Context, queryParam query.QueryParam, deployFilter ListDeploymentQuery) (PaginatedDeployment, *werror.WError) {
	project, werr := s.resourceService.GetResourceByID(ctx, deployFilter.ProjectId)
	if werr != nil {
		return PaginatedDeployment{}, werr
	}

	namespace := deployenv.GetNamespaceName(project.ResourceName, deployFilter.DeploymentEnv)
	deployClient := s.kube.NewDeploymentClient(namespace)
	results, err := deployClient.List(ctx, kubernetes.DeploymentFilter(deployFilter))
	if err != nil {
		return PaginatedDeployment{}, werror.NewFromError(err)
	}

	deployments := []models.Deployment{}
	for _, deploy := range results.Items {
		age := time.Since(deploy.CreationTimestamp.Time)
		condition := deployClient.GetCondition(&deploy)

		status := enum.DeploymentStatusUnhealthy
		if condition.Available.Status == v1.ConditionTrue {
			status = enum.DeploymentStatusHealthy
		}

		deployments = append(deployments, models.Deployment{
			ProjectId:     deployFilter.ProjectId,
			ProjectName:   project.ResourceName,
			DeploymentEnv: deployFilter.DeploymentEnv,
			ServiceName:   strings.TrimSuffix(deploy.Name, "-deployment"),
			Age:           age,
			StringAge:     age.Truncate(time.Second).String(),
			Status:        status,
		})
	}

	deployments = utils.Filter(deployments, func(d models.Deployment) bool {
		return strings.HasPrefix(d.ServiceName, queryParam.QueryFilter.Query)
	})

	slices.SortFunc(deployments, func(a, b models.Deployment) int {
		direction := 0
		switch queryParam.SortFilter.SortBy {
		case "age":
			direction = int(a.Age.Seconds() - b.Age.Seconds())
		case "service_name":
			direction = strings.Compare(a.ServiceName, b.ServiceName)
		case "status":
			direction = strings.Compare(a.Status, b.Status)
		}

		if queryParam.SortFilter.SortOrder == enum.Desc {
			return -direction
		} else {
			return direction
		}
	})

	deployments = utils.Paginate(deployments, queryParam.Pagination.Page, queryParam.Pagination.Limit)

	return PaginatedDeployment{
		Data:  deployments,
		Page:  queryParam.Pagination.Page,
		Limit: queryParam.Pagination.Limit,
		Total: len(deployments),
	}, nil
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
