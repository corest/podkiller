package main

import (
	"log"
	"runtime"
)

func Run(config *killerConfig) error {
    job := killerJob{
		clientset: clientSet(),
		cronstring: config.Scheduler.Crontime,
	}
	scheduler := getJobScheduler(config, &job)
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
