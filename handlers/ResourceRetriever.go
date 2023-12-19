package handlers

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type ResourceInfo struct {
	Kind      string
	Name      string
	Namespace string
	CreatedAt time.Time
}

func ResourceInfos(clientset *kubernetes.Clientset, namespace string) ([]ResourceInfo, error) {
	var resources []ResourceInfo

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, deployment := range deployments.Items {
		resources = append(resources, ResourceInfo{
			Kind:      "Deployment",
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			CreatedAt: deployment.CreationTimestamp.Time,
		})
	}

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, service := range services.Items {
		resources = append(resources, ResourceInfo{
			Kind:      "Service",
			Name:      service.Name,
			Namespace: service.Namespace,
			CreatedAt: service.CreationTimestamp.Time,
		})
	}
	return resources, nil
}

type DeploymentInfo struct {
	Name      string
	Namespace string
	Ready     string
	Age       string
}

func ShowDeployments(clientset *kubernetes.Clientset, namespace string) ([]DeploymentInfo, error) {
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentList []DeploymentInfo
	for _, deployment := range deployments.Items {

		replicaReady := *deployment.Spec.Replicas
		totalReplica := deployment.Status.ReadyReplicas
		deploymentCreatorTimeStamp := deployment.CreationTimestamp
		age := time.Since(deploymentCreatorTimeStamp.Time).Round(time.Second)

		ready := fmt.Sprintf("%v/%v", replicaReady, totalReplica)

		deploymentInfo := DeploymentInfo{
			Name:      deployment.Name,
			Namespace: string(deployment.Namespace),
			Ready:     ready,
			Age:       age.String(),
		}
		deploymentList = append(deploymentList, deploymentInfo)
	}
	return deploymentList, nil
}

type NamespaceInfo struct {
	Name   string
	Status string
	Age    string
}

func NameSpaceShower(clientset *kubernetes.Clientset) ([]NamespaceInfo, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceInfoList []NamespaceInfo
	for _, ns := range namespaces.Items {

		namespaceCreatorTImestamp := ns.GetCreationTimestamp()
		age := time.Since(namespaceCreatorTImestamp.Time).Round(time.Second)

		namespaceInfo := NamespaceInfo{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			Age:    age.String(),
		}
		namespaceInfoList = append(namespaceInfoList, namespaceInfo)
	}
	return namespaceInfoList, nil
}
