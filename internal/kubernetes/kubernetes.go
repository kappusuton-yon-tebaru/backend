package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/kappusuton-yon-tebaru/backend/internal/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	clientset *kubernetes.Clientset
}

func New(cfg *config.Config) (*Kubernetes, error) {
	var config *rest.Config
	var err error
	if cfg.BuilderConfig.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		kubeconfig := filepath.Join(homePath, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Kubernetes{
		clientset,
	}, nil
}
