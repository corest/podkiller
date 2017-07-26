package main

import (
	"log"
	"os"
	"flag"
	"path/filepath"
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
	clientset *kubernetes.Clientset
	cronstring string
}

func (job killerJob) setSchedule(crontstring string) {
    job.cronstring = crontstring
}

func (job killerJob) Run() {
	schedule, _ := cron.Parse(job.cronstring)
    nextrun := schedule.Next(time.Now())
	
	log.Printf("Next pod-killer run at: %s", nextrun.String())

	clientset := job.clientset

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
