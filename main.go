package main

import (
	"log"
	"time"
)

func Run() error {

	for {
		// Calculate duration to sleep before next run
		log.Printf("Print ---")
		time.Sleep(3 * time.Second)
	}
}

func main() {
	// clientSet := clientSet()
	log.Printf("Starting pod-killer...")

	config, err := getConfig()
	if err != nil {
		panic(err.Error())
	}

	log.Printf(config.Title)
	// clientSet := clientSet()

	healthHandler()

	if err := Run(); err != nil {
		panic(err.Error())
	}
}
