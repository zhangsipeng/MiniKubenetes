package cmd

import (
	"example/Minik8s/pkg/apiclient"
	"fmt"

	"github.com/spf13/cobra"
)

var podName_ string = ""
var serviceName_ string = ""
var deploymentName_ string = ""
var serverlessServiceName_ string = ""
var serverlessDAGName_ string = ""
var gpujobsName_ string = ""

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
		var URL string
		if podName_ != "podname" {
			URL = "/api/v1/pods/"
			if podName_ != "all" {
				URL += podName_
			}
		} else if serviceName_ != "servicename" {
			URL = "/api/v1/service/"
			if serviceName_ != "all" {
				URL += serviceName_
			}
		} else if deploymentName_ != "deploymentname" {
			URL = "/api/v1/deployment/"
			if deploymentName_ != "all" {
				URL += deploymentName_
			}
		} else if serverlessServiceName_ != "serverlessServiceName" {
			URL = "/api/v1/serverless/"
			if serverlessServiceName_ != "all" {
				URL += serverlessServiceName_
			}
		} else if serverlessDAGName_ != "serverlessDAGName" {
			URL = "/api/v1/serverlessDAG/"
			if serverlessDAGName_ != "all" {
				URL += serverlessDAGName_
			}
		} else if gpujobsName_ != "gpuName" {
			URL = "/api/v1/gpujobs/"
			if gpujobsName_ != "all" {
				URL += gpujobsName_
			}
		} else {
			fmt.Println("unsupported type")
			return
		}
		// res, err := json.Marshal(msg)
		// if err == nil {
		// 	fmt.Printf("%s msg", res)
		// }
		// reqData := strings.NewReader(string(res))
		// req, _ := http.NewRequest("GET", URL, reqData)
		// req.Header.Add("Content-Type", "application/json")
		// response, err := http.DefaultClient.Do(req)
		// if err != nil && response != nil {
		// 	fmt.Printf("http ret code %d", response.StatusCode)
		// } else {
		// 	fmt.Printf("http err here")
		// }
		fmt.Printf("url %s\n", URL)
		response := apiclient.Request(globalConfig, URL, nil, "GET")
		str_res := string(response[:])
		if str_res == "null" {
			fmt.Printf("workload not found \n")
			return
		}
		// if podName_ == "all" || serviceName_ == "all" || deploymentName_ == "all" {

		// } else {
		// 	var podItem WorkloadResources.Pod
		// 	err := json.Unmarshal(response, &podItem)
		// 	if err != nil {
		// 		panic("Fatel:decode pod json failed \n")
		// 	}
		// 	fmt.Printf("PodName\tStatus\n")
		// 	fmt.Printf("%s\t%s\n", podItem.Metadata.Name, podItem.Status.Phase)
		// }
		fmt.Printf("response is %s\n", response)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringVarP(&podName_, "pod", "p", "podname", "Help message for toggle")
	getCmd.Flags().StringVarP(&serviceName_, "service", "s", "servicename", "Help message for toggle")
	getCmd.Flags().StringVarP(&deploymentName_, "deployment", "d", "deploymentname", "Help message for toggle")
	getCmd.Flags().StringVarP(&serverlessServiceName_, "serverlessService", "", "serverlessServiceName", "Help message for toggle")
	getCmd.Flags().StringVarP(&serverlessDAGName_, "serverlessDAG", "", "serverlessDAGName", "Help message for toggle")
	getCmd.Flags().StringVarP(&gpujobsName_, "gpu", "g", "gpuName", "Help message for toggle")
}
