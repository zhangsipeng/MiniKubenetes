package cmd

import (
	"errors"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/data/WorkloadResources"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func findYamlKind(yamlContent []byte) (kind string, err error) {
	var parsed struct {
		Kind string
	}
	err = yaml.Unmarshal(yamlContent, &parsed)
	if err != nil {
		return
	}
	kind = parsed.Kind
	return
}

var filename string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create -f FILENAME",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create called, filename: %s\n", filename)
		//解析yaml
		if filename == "" {
			panic(errors.New("请指定文件名"))
		}
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("读内容失败", err)
			return
		}
		fmt.Println(string(content))
		kind, err := findYamlKind(content)
		if err != nil {
			fmt.Println("解析失败", err)
		}
		//解析
		switch kind {
		//pod
		case "pod":
			{
				var pod WorkloadResources.Pod
				err = yaml.UnmarshalStrict(content, &pod)
				if err != nil {
					fmt.Println("解析pod失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/pods", pod, "POST")
				fmt.Printf("response is %s\n", response)
				break
			}
		//service
		case "service":
			{
				var service ServiceResources.Service
				err = yaml.UnmarshalStrict(content, &service)
				if err != nil {
					fmt.Println("解析service失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/service", service, "POST")
				fmt.Printf("%s\n", response)
				break
			}
		//replicaset
		case "deployment":
			{
				var deployment WorkloadResources.Deployment
				err = yaml.UnmarshalStrict(content, &deployment)
				if err != nil {
					fmt.Println("解析deployment失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/deployment", deployment, "POST")
				fmt.Printf("%s\n", response)
				break
			}
		// GPU
		case "gpu":
			{
				var job WorkloadResources.GPUJob
				err = yaml.UnmarshalStrict(content, &job)
				if err != nil {
					fmt.Println("解析job失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/gpujobs", job, "POST")
				fmt.Printf("%s\n", response)
				break
			}
		case "serverlessService":
			{
				var serverlessService Serverless.Service
				err = yaml.UnmarshalStrict(content, &serverlessService)
				if err != nil {
					fmt.Println("解析serverless service失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/serverless/", serverlessService, "POST")
				fmt.Printf("%s\n", response)
				break
			}
		case "serverlessDAG":
			{
				var serverlessDAG Serverless.DAG
				err = yaml.UnmarshalStrict(content, &serverlessDAG)
				if err != nil {
					fmt.Println("解析serverless DAG失败", err)
				}
				response := apiclient.Request(globalConfig, "/api/v1/serverlessDAG/", serverlessDAG, "POST")
				fmt.Printf("%s\n", response)
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringVarP(&filename, "file", "f", "", "Help message for toggle")
}
