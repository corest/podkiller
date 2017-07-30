package main

import (
	"log"
	"runtime"
)

func Run(config *killerConfig) error {
	listOptions := getKubernetesListOptions(config)
	job := &killerJob{
		clientset:    clientSet(),
		killerConfig: &config.Killer,
		listOptions:  listOptions,
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
