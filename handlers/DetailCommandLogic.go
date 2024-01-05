package handlers

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodDetails struct {
	Name             string
	Namespace        string
	CreationTime     time.Time
	Phase            string
	Conditions       []PodCondition
	IP               string
	ContainerDetails []ContainerDetails
}

type PodCondition struct {
	Type               string
	Status             string
	LastTransitionTime time.Time
	Reason             string
	Message            string
}

type ContainerDetails struct {
	ContainerName string
	Ports         []PortDetails
}

type PortDetails struct {
	PortName      string
	Protocol      string
	ContainerPort int32
	HostPort      int32
}

func PodDetailsRetrieve(clientset *kubernetes.Clientset, namespace string, podName string) ([]PodDetails, error) {
	poddetail, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var podDetailsList []PodDetails

	pod := poddetail
	podDetails := PodDetails{
		Name:             pod.Name,
		Namespace:        pod.Namespace,
		CreationTime:     pod.CreationTimestamp.Time,
		Phase:            string(pod.Status.Phase),
		Conditions:       make([]PodCondition, len(pod.Status.Conditions)),
		IP:               pod.Status.PodIP,
		ContainerDetails: make([]ContainerDetails, len(pod.Spec.Containers)),
	}

	for i, condition := range pod.Status.Conditions {
		podDetails.Conditions[i] = PodCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime.Time,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	for i, container := range pod.Spec.Containers {
		containerDetails := ContainerDetails{
			ContainerName: container.Name,
			Ports:         make([]PortDetails, len(container.Ports)),
		}
		for j, port := range container.Ports {
			containerDetails.Ports[j] = PortDetails{
				PortName:      port.Name,
				Protocol:      string(port.Protocol),
				ContainerPort: port.ContainerPort,
				HostPort:      port.HostPort,
			}
		}
		podDetails.ContainerDetails[i] = containerDetails
	}

	podDetailsList = append(podDetailsList, podDetails)

	return podDetailsList, nil
}

type DeploymentDetails struct {
	Name              string
	Namespace         string
	CreationTime      time.Time
	Replicas          int32
	AvailableReplicas int32
	ReadyReplicas     int32
	UpdatedReplicas   int32
	Strategy          string
	Selector          string
	Containers        []ContainerDetails
}

func GetDeploymentDetails(clientset *kubernetes.Clientset, namespace string, deploymentName string) ([]DeploymentDetails, error) {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentDetailsList []DeploymentDetails

	deploy := deployment
	deploymentDetails := DeploymentDetails{
		Name:              deploy.Name,
		Namespace:         deploy.Namespace,
		CreationTime:      deploy.CreationTimestamp.Time,
		Replicas:          *deploy.Spec.Replicas,
		AvailableReplicas: deploy.Status.AvailableReplicas,
		ReadyReplicas:     deploy.Status.ReadyReplicas,
		UpdatedReplicas:   deploy.Status.UpdatedReplicas,
		Strategy:          string(deploy.Spec.Strategy.Type),
		Selector:          getLabelSelector(deploy.Spec.Selector),
		Containers:        make([]ContainerDetails, len(deploy.Spec.Template.Spec.Containers)),
	}

	for i, container := range deploy.Spec.Template.Spec.Containers {
		containerDetails := ContainerDetails{
			ContainerName: container.Name,
			Ports:         make([]PortDetails, len(container.Ports)),
		}
		for j, port := range container.Ports {
			containerDetails.Ports[j] = PortDetails{
				PortName:      port.Name,
				Protocol:      string(port.Protocol),
				ContainerPort: port.ContainerPort,
			}
		}
		deploymentDetails.Containers[i] = containerDetails
	}

	deploymentDetailsList = append(deploymentDetailsList, deploymentDetails)

	return deploymentDetailsList, nil
}

func getLabelSelector(selector *v1.LabelSelector) string {
	if selector == nil || len(selector.MatchLabels) == 0 {
		return ""
	}

	labelSelector := ""
	for key, value := range selector.MatchLabels {
		labelSelector += key + "=" + value + ","
	}

	// Remove the trailing comma
	labelSelector = labelSelector[:len(labelSelector)-1]

	return labelSelector
}

type NamespaceDetails struct {
	Name          string
	CreationTime  time.Time
	Status        string
	Labels        map[string]string
	Annotations   map[string]string
	ResourceQuota string
}

func NameSpaceDetailsRetrieve(clientset *kubernetes.Clientset, namespace string) ([]NamespaceDetails, error) {
	ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var nsDetailsList []NamespaceDetails

	nsDetails := NamespaceDetails{
		Name:          ns.Name,
		CreationTime:  ns.CreationTimestamp.Time,
		Status:        string(ns.Status.Phase),
		Labels:        ns.Labels,
		Annotations:   ns.Annotations,
		ResourceQuota: getResourceQuota(clientset, namespace),
	}

	nsDetailsList = append(nsDetailsList, nsDetails)

	return nsDetailsList, nil
}

func getResourceQuota(clientset *kubernetes.Clientset, namespace string) string {
	quota, err := clientset.CoreV1().ResourceQuotas(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return ""
	}

	if len(quota.Items) > 0 {
		return quota.Items[0].Name
	}

	return ""
}

type ServiceDetails struct {
	Name            string
	Namespace       string
	CreationTime    time.Time
	Labels          map[string]string
	Type            string
	ClusterIP       string
	ExternalIPs     []string
	LoadBalancerIP  string
	Ports           []ServicePortDetails
	Selector        map[string]string
	SessionAffinity string
}

type ServicePortDetails struct {
	Name       string
	Protocol   string
	Port       int32
	TargetPort string
	NodePort   int32
}

func ServiceDetailsRetrieve(clientset *kubernetes.Clientset, namespace string, serviceName string) ([]ServiceDetails, error) {
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var serviceDetailsList []ServiceDetails

	serviceDetails := ServiceDetails{
		Name:            service.Name,
		Namespace:       service.Namespace,
		CreationTime:    service.CreationTimestamp.Time,
		Labels:          service.Labels,
		Type:            string(service.Spec.Type),
		ClusterIP:       service.Spec.ClusterIP,
		ExternalIPs:     service.Spec.ExternalIPs,
		LoadBalancerIP:  service.Spec.LoadBalancerIP,
		Ports:           make([]ServicePortDetails, len(service.Spec.Ports)),
		Selector:        service.Spec.Selector,
		SessionAffinity: string(service.Spec.SessionAffinity),
	}

	for i, port := range service.Spec.Ports {
		servicePortDetails := ServicePortDetails{
			Name:       port.Name,
			Protocol:   string(port.Protocol),
			Port:       port.Port,
			TargetPort: port.TargetPort.String(),
			NodePort:   port.NodePort,
		}
		serviceDetails.Ports[i] = servicePortDetails
	}

	serviceDetailsList = append(serviceDetailsList, serviceDetails)

	return serviceDetailsList, nil
}
