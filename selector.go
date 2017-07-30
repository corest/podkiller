package main

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func sliceOperation(op string, slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m1 := map[string]int{}
	m2 := map[string]int{}

	for _, v := range slice1 {
		m1[v] = 1
	}
	for _, v := range slice2 {
		m2[v] = 1
	}

	switch op {
	case "substruction":
		for k, _ := range m1 {
			if _, contains := m2[k]; !contains {
				diffStr = append(diffStr, k)
			}
		}
	case "unity":
		for k, _ := range m1 {
			if _, contains := m2[k]; contains {
				diffStr = append(diffStr, k)
			}
		}
	default:
		log.Fatalf("Unsupported operation for sliceOperation %s", op)
	}

	return diffStr
}

func getKubernetesListOptions(config *killerConfig) *metav1.ListOptions {

	if !config.Killer.Match_all {
		if reqs, err := labels.ParseToRequirements(config.Killer.Selector); err != nil {
			log.Fatalf("Failed to create requirement from reqs %s \n %v", reqs, err)
		}

	}

	return &metav1.ListOptions{LabelSelector: config.Killer.Selector}
}

func getKubernetesNamespaces(config *runnerConfig, clientset *kubernetes.Clientset) []string {
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	var existingNamespaces []string
	var resultingNamespaces []string
	if err != nil {
		log.Fatalf("Failed to get list of namespaces %v", err)
	}
	for _, namespace := range namespaces.Items {
		existingNamespaces = append(existingNamespaces, namespace.Name)
	}
	if config.Namespace_deny_policy {
		resultingNamespaces = sliceOperation("substruction", existingNamespaces, config.Namespace_list)
	} else {
		resultingNamespaces = sliceOperation("unity", existingNamespaces, config.Namespace_list)
	}

	return resultingNamespaces
}
