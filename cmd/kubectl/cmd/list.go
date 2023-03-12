package cmd

import (
	"example/Minik8s/pkg/apiclient"
	"fmt"

	"github.com/spf13/cobra"
)

var listType string = ""

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		if listType != "pods" && listType != "nodes" && listType != "service" && listType != "gpujobs" &&
			listType != "deployment" && listType != "serverless" && listType != "serverlessDAG" {
			fmt.Println("unsupported type")
			return
		}
		URL := "/api/v1/" + listType

		fmt.Printf("url %s\n", URL)
		response := apiclient.Request(globalConfig, URL, nil, "GET")
		str_res := string(response[:])
		if str_res == "null" {
			fmt.Printf("workload not found \n")
			return
		}

		fmt.Printf("response is %s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().StringVarP(&listType, "type", "t", "type", "Help message for toggle")
}
