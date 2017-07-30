package main

import (
	"log"
	"runtime"
)

func Run(config *killerConfig) error {
	clientset := clientSet()
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
		log.Fatalf("Unable to schedule pod-killer", err)
	}

	scheduler.Start()
	return nil
}

func main() {
	// clientSet := clientSet()
	log.Printf("Starting pod-killer...")

	config, err := getConfig()
	if err != nil {
		panic(err.Error())
	}

	healthHandler(&config)

	if err := Run(&config); err != nil {
		panic(err.Error())
	}

	runtime.Goexit()
}
