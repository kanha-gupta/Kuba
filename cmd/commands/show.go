/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"github.com/kanha-gupta/kuba/cmd"
	"github.com/kanha-gupta/kuba/handlers"
	"github.com/kanha-gupta/kuba/kubernetesClient"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var Verbose bool
var Source string
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "It show the rosurce details metion in the command (like.. pods, servcies, deployments, etc...)",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Please mention the name of resorces you want to see")
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all resources from provided namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")
		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes k8_client: %v", err)
		}
		resources, err := handlers.ResourceInfos(client, namespace)
		if err != nil {
			log.Printf("error getting resources: %v", err)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Resource Type", "Name", "Namespace", "Created At"})

			for _, resource := range resources {
				createdTime := resource.CreatedAt.Format("2006-01-02 15:04:05")
				row := []string{resource.Kind, resource.Name, resource.Namespace, createdTime}
				table.Append(row)
			}
			table.Render()
		}
	},
}

var deploymentCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Show deployments in a Kubernetes namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("ns")
		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes k8_client: %v", err)
		}
		deploymentList, err := handlers.ShowDeployments(client, namespace)
		if err != nil {
			log.Printf("error getting deployment list: %v", err)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Deployment", "Namespace", "Ready", "Age"})

			for _, deployment := range deploymentList {
				row := []string{deployment.Name, deployment.Namespace, deployment.Ready, deployment.Age}
				table.Append(row)
			}
			table.Render()
		}
	},
}

var namespaceCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "It will show all name-spaces in kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubernetesClient.GetClient()
		if err != nil {
			log.Printf("error getting kubernetes k8_client: %v", err)
		}
		namespaceDetails, err := handlers.NameSpaceShower(client)
		if err != nil {
			log.Printf("Can't get the namespacces: %v", err)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Namespace-Name", "status", "Age"})

			for _, namespace := range namespaceDetails {
				row := []string{namespace.Name, namespace.Status, namespace.Age}
				table.Append(row)
			}
			table.Render()
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(showCmd)
	showCmd.AddCommand(allCmd)
	showCmd.AddCommand(namespaceCmd)
	showCmd.AddCommand(deploymentCmd)
}
