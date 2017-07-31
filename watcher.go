package main

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb/client/v2"

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

	var clnt client.Client

	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", "localhost", 8086),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  podkillerDb,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	var fields map[string]interface{}
	tags := map[string]string{"service": podkillerDbTag}

	fields = map[string]interface{}{
		"name":          pod.ObjectMeta.Name,
		"namespace":     pod.ObjectMeta.Namespace,
		"condemnded_at": time.Now(),
	}
	pt, err := client.NewPoint(podkillerDbSeries, tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := clnt.Write(bp); err != nil {
		log.Fatal(err)
	}

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
