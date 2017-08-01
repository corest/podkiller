package main

import (
	"log"
	"runtime"

	"k8s.io/client-go/kubernetes"
)

func run(config *Config, clientset *kubernetes.Clientset,
	influxmanager *InfluxManager) error {

	listOptions := getKubernetesListOptions(config)
	allowedNamespaces := getKubernetesNamespaces(config, clientset)

	job := &KillerJob{
		clientset:         clientset,
		config:            &config.Killer,
		influxmanager:     influxmanager,
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

	influxcclient := getInfluxClient(&config)
	influxmanager := InfluxManager{client: influxcclient}
	influxmanager.initDB()

	routesHandler(&config, &influxmanager)

	watchPods(clientset)

	if err := run(&config, clientset, &influxmanager); err != nil {
		panic(err.Error())
	}

	runtime.Goexit()
}
