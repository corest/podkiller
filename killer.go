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
	"k8s.io/apimachinery/pkg/api/errors"
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
	clientset    *kubernetes.Clientset
	killerConfig *runnerConfig
	cronstring   string
}

type eventReceivers struct {
	pods []*doomedPod
}

func (cemetry *eventReceivers) layPodLow(pod *doomedPod) error {
	cemetry.pods = append(cemetry.pods, pod)
	return nil
}

type doomedPod struct {
	name    string
	isAlive bool
}

func (job *killerJob) setSchedule(crontstring string) {
	job.cronstring = crontstring
}

func (job *killerJob) executeDoomedPod(pods chan *doomedPod, necrolog *chan *doomedPod, wg *sync.WaitGroup) error {
	defer wg.Done()
	var pod *doomedPod
	select {
	case pod = <-pods:
		pod.isAlive = false
		log.Printf("Executing pod '%s'", pod.name)
	default:
		fmt.Println("No pods were executed")
	}

	select {
	case *necrolog <- pod:
		log.Printf("Mark pod '%s' as dead ", pod.name)
	default:
		fmt.Println("No pods were marked as dead")
	}

	return nil
}

func (job killerJob) Run() {

	var wg sync.WaitGroup

	// get somewhere dooomed pods by puutting into channel and return channel https://blog.golang.org/pipelines
	testPod1 := &doomedPod{
		name:    "Testpod1",
		isAlive: true,
	}
	testPod2 := &doomedPod{
		name:    "Testpod2",
		isAlive: true,
	}

	num := 2

	wg.Add(num)
	condemnedPods := make(chan *doomedPod, num)
	necrology := make(chan *doomedPod, num) // buffered channel because unbeffered wait for read

	condemnedPods <- testPod1
	condemnedPods <- testPod2

	for i := 0; i < num; i++ {
		go job.executeDoomedPod(condemnedPods, &necrology, &wg)
	}

	go func() {
		for pod := range necrology {
			// use this for sending events
			log.Printf("Pod '%s' is alive %t", pod.name, pod.isAlive)
		}
	}()

	wg.Wait()
	/* 	clientset := job.clientset

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	_, err = clientset.CoreV1().Pods("default").Get("nginx", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Pod not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		log.Printf("Found pod\n")
	} */
	schedule, _ := cron.Parse(job.cronstring)
	nextrun := schedule.Next(time.Now())
	log.Printf("The Moor has done his work, the Moor can go. Next run at: %s", nextrun.String())
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

func listPods() {
	clientset := clientSet()
	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		log.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get("nginx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			log.Printf("Pod not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			log.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			log.Printf("Found pod\n")
		}

		time.Sleep(10 * time.Second)
	}
}

func killPod(clientset *kubernetes.Clientset, namespace string, pod string, reason string) {
	log.Printf("Killing pod %s because %s\n", pod, reason)
	err := clientset.Core().Pods(namespace).Delete(pod, nil)
	if err != nil {
		log.Printf("Unable to delete pod %s because %s", pod, err)
	}
}
