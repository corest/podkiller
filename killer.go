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

type killerJob struct {
	clientset         *kubernetes.Clientset
	killerConfig      *runnerConfig
	listOptions       *metav1.ListOptions
	allowedNamespaces []string
	cronstring        string
}

type namespaceNecrology struct {
	namespace string
	pods      []*doomedPod
}

func (cemetry *namespaceNecrology) layPodLow(pod *doomedPod) error {
	cemetry.pods = append(cemetry.pods, pod)
	return nil
}

type doomedPod struct {
	name      string
	namespace string
	isAlive   bool
}

func (job *killerJob) setSchedule(crontstring string) {
	job.cronstring = crontstring
}

func (job *killerJob) extractDoomedPods(namespaces []string) ([]*doomedPod, error) {
	var condemnedPods []*doomedPod
	for _, namespace := range namespaces {
		pods, err := job.clientset.CoreV1().Pods(namespace).List(*job.listOptions)
		if err != nil {
			log.Printf("Failed to get pods from namespace '%s'. %v ", namespace, err)
		}
		for _, pod := range pods.Items {
			dpod := &doomedPod{
				name:      pod.Name,
				namespace: namespace,
				isAlive:   true,
			}
			condemnedPods = append(condemnedPods, dpod)
		}
	}
	return condemnedPods, nil
}

func (job *killerJob) executeDoomedPods(namespace string, wg *sync.WaitGroup) error {
	defer wg.Done()
	pods, err := job.clientset.CoreV1().Pods(namespace).List(*job.listOptions)
	if err != nil {
		log.Printf("Failed to get pods from namespace '%s'. %v ", namespace, err)
	}

	if len(pods.Items) > 0 {
		podsNumber := len(pods.Items)
		log.Printf("Found %d pods to kill in namespace '%s'.", podsNumber, namespace)

		condemnedPodsChannel := make(chan *doomedPod, podsNumber)
		for _, pod := range pods.Items {
			dpod := &doomedPod{
				name:      pod.Name,
				namespace: namespace,
				isAlive:   true,
			}
			condemnedPodsChannel <- dpod
		}

		necrology := make(chan *doomedPod, podsNumber)
		var podsWG sync.WaitGroup
		podsWG.Add(podsNumber)
		for i := 0; i < len(pods.Items); i++ {
			go func() {
				defer podsWG.Done()
				var pod *doomedPod
				select {
				case pod = <-condemnedPodsChannel:
					pod.isAlive = false
					log.Printf("Executing pod '%s'(namespace: '%s')", pod.name, pod.namespace)
					err := job.clientset.Core().Pods(pod.namespace).Delete(pod.name, nil)
					if err != nil {
						log.Printf("Unable to delete pod %s(namespace: '%s') because %v", pod.name, pod.namespace, err)
					}
				default:
					fmt.Println("No pods were executed")
				}

				select {
				case necrology <- pod:
					log.Printf("Mark pod '%s'(namespace: '%s') as dead ", pod.name, pod.namespace)
				default:
					fmt.Println("No pods were marked as dead")
				}
			}()
		}

		podsWG.Wait()
	}

	return nil
}

func (job killerJob) Run() {

	log.Println("Running new killer job...")

	var namespaceWg sync.WaitGroup

	namespaces := job.allowedNamespaces
	namespaceNumber := len(namespaces)
	namespaceWg.Add(namespaceNumber)

	for _, namespace := range namespaces {
		go job.executeDoomedPods(namespace, &namespaceWg)
	}

	namespaceWg.Wait()

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
