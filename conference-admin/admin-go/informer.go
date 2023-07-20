package main

import (
	"time"

	"github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/api/types/v1alpha1"
	client_v1alpha1 "github.com/salaboy/platforms-on-k8s/conference-admin/admin-go/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func WatchResources(clientSet client_v1alpha1.ConferenceAdminV1Alpha1Interface) cache.Store {
	environmentStore, environmentController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.Environments("default").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Environments("default").Watch(lo)
			},
		},
		&v1alpha1.Environment{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go environmentController.Run(wait.NeverStop)
	return environmentStore
}
