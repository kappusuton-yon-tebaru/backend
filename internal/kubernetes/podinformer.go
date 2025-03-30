package kubernetes

import (
	"time"

	"k8s.io/api/core/v1"
	apicorev1 "k8s.io/api/core/v1"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type PodInformerEventHandler interface {
	PodCreated(pod *v1.Pod)
	PodUpdated(oldPod *v1.Pod, newPod *v1.Pod)
	PodDeleted(pod *v1.Pod)
}

type PodInformer struct {
	stopper chan struct{}
}

func (kube *Kubernetes) NewPodInformer(handlers PodInformerEventHandler, podSelector string) (*PodInformer, error) {
	factory := informers.NewSharedInformerFactoryWithOptions(kube.clientset, time.Minute, informers.WithTweakListOptions(
		func(opts *apimetav1.ListOptions) {
			opts.LabelSelector = podSelector
		}),
	)

	podInformer := factory.Core().V1().Pods().Informer()
	_, err := podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			pod := obj.(*apicorev1.Pod)
			handlers.PodCreated(pod)
		},
		UpdateFunc: func(oldObj, newObj any) {
			oldPod := oldObj.(*apicorev1.Pod)
			newPod := newObj.(*apicorev1.Pod)
			handlers.PodUpdated(oldPod, newPod)
		},
		DeleteFunc: func(obj any) {
			pod := obj.(*apicorev1.Pod)
			handlers.PodDeleted(pod)
		},
	})

	if err != nil {
		return nil, err
	}

	stopper := make(chan struct{})

	factory.Start(stopper)
	factory.WaitForCacheSync(stopper)

	return &PodInformer{
		stopper,
	}, nil
}

func (p *PodInformer) Stop() {
	close(p.stopper)
}
