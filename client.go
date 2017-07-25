package main

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func clientSet() *kubernetes.Clientset {
	log.Printf("Setup kubernetes client")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet
}

func reap(clientSet *kubernetes.Clientset, namespace string, pod string, reason string) {
	fmt.Printf("Reaping Pod %s because %s\n", pod, reason)
	err := clientSet.Core().Pods(namespace).Delete(pod, nil)
	if err != nil {
		// log the error, but continue on
		fmt.Fprintf(os.Stderr, "unable to delete pod %s because %s", pod, err)
	}
}
