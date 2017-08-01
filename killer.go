package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/robfig/cron"

	/*	"k8s.io/client-go/kubernetes"
		"k8s.io/client-go/rest"*/

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*func clientset() *kubernetes.Clientset {
	log.Printf("Setup kubernetes client")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}*/

// KillerJob structure
type KillerJob struct {
	clientset         *kubernetes.Clientset
	config            *KillerConfig
	listOptions       *metav1.ListOptions
	influxmanager     *InfluxManager
	allowedNamespaces []string
	cronstring        string
	scheduletAt       time.Time
}

// DoomedPod structure
type DoomedPod struct {
	name        string
	namespace   string
	isAlive     bool
	condemnedAt time.Time
}

func (job *KillerJob) setSchedule(crontstring string) {
	job.cronstring = crontstring
}

func (job *KillerJob) executeDoomedPods(namespace string, wg *sync.WaitGroup) error {
	defer wg.Done()
	pods, err := job.clientset.CoreV1().Pods(namespace).List(*job.listOptions)
	if err != nil {
		log.Printf("Failed to get pods from namespace '%s'. %v ", namespace, err)
	}

	if len(pods.Items) > 0 {
		podsNumber := len(pods.Items)
		log.Printf("Found %d pods to kill in namespace '%s'.", podsNumber, namespace)

		condemnedPodsChannel := make(chan *DoomedPod, podsNumber)
		for _, pod := range pods.Items {
			dpod := &DoomedPod{
				name:        pod.Name,
				namespace:   namespace,
				isAlive:     true,
				condemnedAt: job.scheduletAt,
			}
			condemnedPodsChannel <- dpod
		}

		var podsWG sync.WaitGroup
		podsWG.Add(podsNumber)
		for i := 0; i < len(pods.Items); i++ {
			go func() {
				defer podsWG.Done()
				var pod *DoomedPod
				select {
				case pod = <-condemnedPodsChannel:
					pod.isAlive = false
					log.Printf("Killing pod '%s'(namespace: '%s')...", pod.name, pod.namespace)
					err := job.clientset.Core().Pods(pod.namespace).Delete(pod.name, nil)
					if err != nil {
						log.Printf("Unable to delete pod %s(namespace: '%s') because %v", pod.name, pod.namespace, err)
					}
					job.influxmanager.addDoomedPod(pod)

				default:
					fmt.Println("No pods were executed")
				}
			}()
		}

		podsWG.Wait()
	}

	return nil
}

// Run - main function for running scheduled jobs
func (job KillerJob) Run() {

	log.Println("Running new killer job...")

	var namespaceWg sync.WaitGroup

	job.scheduletAt = time.Now()
	namespaces := job.allowedNamespaces
	namespaceNumber := len(namespaces)
	namespaceWg.Add(namespaceNumber)

	for _, namespace := range namespaces {
		go job.executeDoomedPods(namespace, &namespaceWg)
	}

	namespaceWg.Wait()

	if err := job.influxmanager.writePodNecrology(); err != nil {
		log.Fatalf("Failed to write new point into influxdb %v", err)
	}

	schedule, _ := cron.Parse(job.cronstring)
	nextrun := schedule.Next(time.Now())
	log.Printf("...done. Next run at: %s", nextrun.String())
}

func clientSet() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
