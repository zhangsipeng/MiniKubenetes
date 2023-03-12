package cmd

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/Serverless"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// updateCmd represents the create command
var updateCmd = &cobra.Command{
	Use:   "update -f FILENAME",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create called ,filename : %s", filename)
		//解析yaml
		if filename == "" {
			fmt.Println("请指定一个文件名")
			return
		}
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("读文件失败", err)
			return
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("读内容失败", err)
			return
		}
		fmt.Println(string(content))
		err = file.Close()
		if err != nil {
			fmt.Println("关文件失败", err)
			return
		}
		kind, err := findYamlKind(content)
		if err != nil {
			fmt.Println("解析失败", err)
		}
		//解析
		switch kind {
		//pod
		case "serverlessService":
			{
				var serverlessService Serverless.Service
				err = yaml.UnmarshalStrict(content, &serverlessService)
				if err != nil {
					fmt.Println("解析serverlessService失败", err)
				}
				serverlessService.Status.Phase = "change"
				response := apiclient.Request(globalConfig, "/api/v1/serverless/", serverlessService, "PUT")
				fmt.Printf("response is %s\n", response)
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateCmd.Flags().StringVarP(&filename, "file", "f", "", "Help message for toggle")
}
