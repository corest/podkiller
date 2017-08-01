package main

import (
	"fmt"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func sliceOperation(op string, slice1 []string, slice2 []string) ([]string, error) {
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
		for k := range m1 {
			if _, contains := m2[k]; !contains {
				diffStr = append(diffStr, k)
			}
		}
	case "unity":
		for k := range m1 {
			if _, contains := m2[k]; contains {
				diffStr = append(diffStr, k)
			}
		}
	default:
		return []string{}, fmt.Errorf("Unsupported operation for sliceOperation %s", op)
	}

	return diffStr, nil
}

func getKubernetesListOptions(config *Config) *metav1.ListOptions {

	if reqs, err := labels.ParseToRequirements(config.Killer.Selector); err != nil {
		log.Fatalf("Failed to create requirement from reqs %s \n %v", reqs, err)
	}

	log.Printf("Used selector: '%s'", config.Killer.Selector)

	return &metav1.ListOptions{LabelSelector: config.Killer.Selector}
}

func getKubernetesNamespaces(config *Config, clientset *kubernetes.Clientset) []string {
	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	var existingNamespaces []string
	var resultingNamespaces []string
	if err != nil {
		log.Fatalf("Failed to get list of namespaces %v", err)
	}
	for _, namespace := range namespaces.Items {
		existingNamespaces = append(existingNamespaces, namespace.Name)
	}
	if config.Killer.NamespaceDenyPolicy {
		resultingNamespaces, err = sliceOperation("substruction", existingNamespaces, config.Killer.NamespaceList)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		resultingNamespaces, err = sliceOperation("unity", existingNamespaces, config.Killer.NamespaceList)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	log.Printf("Allowed namespaces for executing pod kills on: [%s]",
		strings.Join(resultingNamespaces, ", "))

	return resultingNamespaces
}
