package cmd

import (
	"example/Minik8s/pkg/apiclient"
	"fmt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var podName string = ""
var serviceName string = ""
var deploymentName string = ""
var serverlessServiceName string = ""
var deleteCmd = &cobra.Command{
	Use:   "delete ",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// var URL string
		if podName != "podname" {
			response := apiclient.Request(globalConfig, "/api/v1/pods/"+podName, nil, "DELETE")
			fmt.Printf("response is %s", response)
		} else if serviceName != "servicename" {
			response := apiclient.Request(globalConfig, "/api/v1/service/"+serviceName, nil, "DELETE")
			fmt.Printf("response is %s", response)
		} else if deploymentName != "deploymentname" {
			response := apiclient.Request(globalConfig, "/api/v1/deployment/"+deploymentName, nil, "DELETE")
			fmt.Printf("response is %s", response)
		} else {
			response := apiclient.Request(globalConfig, "/api/v1/serverless/"+serverlessServiceName, nil, "DELETE")
			fmt.Printf("response is %s", response)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().StringVarP(&podName, "pod", "p", "podname", "Help message for toggle")
	deleteCmd.Flags().StringVarP(&serviceName, "service", "s", "servicename", "Help message for toggle")
	deleteCmd.Flags().StringVarP(&deploymentName, "deployment", "d", "deploymentname", "Help message for toggle")
	deleteCmd.Flags().StringVarP(&serverlessServiceName, "serverlessService", "", "serverlessServiceName", "Help message for toggle")
}
