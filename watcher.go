package main

import (
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func podDeleted(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Printf("Event: pod '%s' from namespace '%s' deleted at %s",
		pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, time.Now())
}

func watchPods(client *kubernetes.Clientset) (cache.Store, error) {
	watchlist := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods",
		v1.NamespaceAll, fields.Everything())
	resyncPeriod := 30 * time.Minute
	store, controller := cache.NewInformer(watchlist, &v1.Pod{}, resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: podDeleted,
		})
	go controller.Run(wait.NeverStop)
	return store, nil
}
