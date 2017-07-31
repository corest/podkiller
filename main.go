package main

import (
	"log"
	"runtime"

	"k8s.io/client-go/kubernetes"
)

func run(config *killerConfig, clientset *kubernetes.Clientset) error {
	listOptions := getKubernetesListOptions(config)
	allowedNamespaces := getKubernetesNamespaces(config, clientset)
	job := &killerJob{
		clientset:         clientset,
		killerConfig:      &config.Killer,
		listOptions:       listOptions,
		allowedNamespaces: allowedNamespaces,
	}
	scheduler, err := getJobScheduler(config, job)
	if err != nil {
		log.Fatalf("Unable to schedule pod-killer %v", err)
	}

	scheduler.Start()
	return nil
}

func main() {
	log.Printf("Starting pod-killer...")

	config, err := getConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset := clientSet()

	healthHandler(&config)

	watchPods(clientset)

	if err := run(&config, clientset); err != nil {
		panic(err.Error())
	}

	runtime.Goexit()
}
