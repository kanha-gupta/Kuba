package commands

import (
	"fmt"
	"github.com/kanha-gupta/kuba/cmd"
	"github.com/kanha-gupta/kuba/handlers"
	"github.com/kanha-gupta/kuba/kubernetesClient"
	"github.com/spf13/cobra"
	"log"
)

var DetailsCommand = &cobra.Command{
	Use:   "details",
	Short: "show kubernetes resources details",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("please provide resource type to get details")
	},
}

var podCommand = &cobra.Command{
	Use:   "pod",
	Short: "Show details of a Kubernetes pod",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")
		podName, _ := cmd.Flags().GetString("p")

		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes client: %v", err)
		}
		podDetailsList, err := handlers.PodDetailsRetrieve(client, namespace, podName)
		if err != nil {
			log.Printf("error getting pod details: %v", err)
		} else {
			for _, pod := range podDetailsList {
				fmt.Println("Name:", pod.Name)
				fmt.Println("Namespace:", pod.Namespace)
				fmt.Println("Creation Time:", pod.CreationTime)
				fmt.Println("Phase:", pod.Phase)
				fmt.Println("IP:", pod.IP)

				fmt.Println("Conditions:")
				for _, condition := range pod.Conditions {
					fmt.Println("\tType:", condition.Type)
					fmt.Println("\tStatus:", condition.Status)
					fmt.Println("\tLast Transition Time:", condition.LastTransitionTime)
					fmt.Println("\tReason:", condition.Reason)
					fmt.Println("\tMessage:", condition.Message)
				}

				fmt.Println("Container Details:")
				for _, container := range pod.ContainerDetails {
					fmt.Println("\tContainer Name:", container.ContainerName)
					fmt.Println("\tPorts:")
					for _, port := range container.Ports {
						fmt.Println("\t\tPort Name:", port.PortName)
						fmt.Println("\t\tProtocol:", port.Protocol)
						fmt.Println("\t\tContainer Port:", port.ContainerPort)
						fmt.Println("\t\tHost Port:", port.HostPort)
					}
				}

				fmt.Println("-----------------------------------")
			}

		}
	},
}

var deploymentCommand = &cobra.Command{
	Use:   "deployment",
	Short: "Show details of a Kubernetes deployment",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")
		deploymentName, _ := cmd.Flags().GetString("d")

		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes client: %v", err)
		}
		deploymentDetailsList, err := handlers.GetDeploymentDetails(client, namespace, deploymentName)
		if err != nil {
			log.Printf("error getting deployment details: %v", err)
		} else {
			for _, deployment := range deploymentDetailsList {
				fmt.Println("Name:", deployment.Name)
				fmt.Println("Namespace:", deployment.Namespace)
				fmt.Println("Creation Time:", deployment.CreationTime)
				fmt.Println("Replicas:", deployment.Replicas)
				fmt.Println("Available Replicas:", deployment.AvailableReplicas)
				fmt.Println("Ready Replicas:", deployment.ReadyReplicas)
				fmt.Println("Updated Replicas:", deployment.UpdatedReplicas)
				fmt.Println("Strategy:", deployment.Strategy)
				fmt.Println("Selector:", deployment.Selector)

				fmt.Println("Containers:")
				for _, container := range deployment.Containers {
					fmt.Println("\tContainer Name:", container.ContainerName)
					fmt.Println("\tPorts:")
					for _, port := range container.Ports {
						fmt.Println("\t\tPort Name:", port.PortName)
						fmt.Println("\t\tProtocol:", port.Protocol)
						fmt.Println("\t\tContainer Port:", port.ContainerPort)
					}
				}

				fmt.Println("-----------------------------------")
			}

		}
	},
}
var namespaceCommand = &cobra.Command{
	Use:   "namespace",
	Short: "Show details of a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")

		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes client: %v", err)
			return
		}

		nsDetailsList, err := handlers.NameSpaceDetailsRetrieve(client, namespace)
		if err != nil {
			log.Printf("error getting namespace details: %v", err)
			return
		}

		for _, ns := range nsDetailsList {
			fmt.Println("Name:", ns.Name)
			fmt.Println("Creation Time:", ns.CreationTime)
			fmt.Println("Status:", ns.Status)
			fmt.Println("Labels:", ns.Labels)
			fmt.Println("Annotations:", ns.Annotations)
			fmt.Println("Resource Quota:", ns.ResourceQuota)
			fmt.Println("-----------------------------------")
		}
	},
}
var serviceCommand = &cobra.Command{
	Use:   "service",
	Short: "Get details of a Kubernetes service",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")
		serviceName, _ := cmd.Flags().GetString("s")

		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes client: %v", err)
		}
		serviceDetailsList, err := handlers.ServiceDetailsRetrieve(client, namespace, serviceName)
		if err != nil {
			log.Printf("error getting service details: %v", err)
		} else {
			for _, service := range serviceDetailsList {
				fmt.Println("Service Name:", service.Name)
				fmt.Println("Namespace:", service.Namespace)
				fmt.Println("Creation Time:", service.CreationTime)
				fmt.Println("Labels:", service.Labels)
				fmt.Println("Type:", service.Type)
				fmt.Println("Cluster IP:", service.ClusterIP)
				fmt.Println("External IPs:", service.ExternalIPs)
				fmt.Println("LoadBalancer IP:", service.LoadBalancerIP)

				fmt.Println("Ports:")
				for _, port := range service.Ports {
					fmt.Println("  - Name:", port.Name)
					fmt.Println("    Protocol:", port.Protocol)
					fmt.Println("    Port:", port.Port)
					fmt.Println("    Target Port:", port.TargetPort)
					fmt.Println("    Node Port:", port.NodePort)
				}

				fmt.Println("Selector:", service.Selector)
				fmt.Println("Session Affinity:", service.SessionAffinity)

				fmt.Println("---------------------------")
			}
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(DetailsCommand)
	DetailsCommand.AddCommand(podCommand)
	podCommand.PersistentFlags().String("p", "", "You need to provide the name of pod in order to get details of that perticular pod (eg: --p=pod-name)")
	DetailsCommand.AddCommand(deploymentCommand)
	deploymentCommand.PersistentFlags().String("d", "", "You need to provide the name of deployment to get details (eg: --d=deployment-name)")
	DetailsCommand.AddCommand(namespaceCommand)
	namespaceCommand.PersistentFlags().String("ns", "", "Provide the name of the namespace to get its details (e.g., --ns=namespace-name)")
	DetailsCommand.AddCommand(serviceCommand)
	serviceCommand.PersistentFlags().String("s", "", "You need to provide the name of pod in order to get details of that perticular pod (eg: --s=service-name)")
}
