package main

import (
	"example/Minik8s/cmd/kubectl/cmd"
	"example/Minik8s/pkg/apiclient"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type authConfig struct {
	Init           bool
	ApiServerIP    string
	Token          string
	CaHash         string
	ConfigFileName string
}

func main() {
	// TODO:change it to read file from local fs
	// info := apiclient.GetInitInfo()
	// runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	ConfigFileName := "authconfig.yaml"
	absPath, _ := filepath.Abs(ConfigFileName)
	file, err := os.Open(absPath)
	if err != nil {
		fmt.Println("配置文件打开失败", err)
		return
	}
	content, err := ioutil.ReadAll(file)
	fmt.Println(string(content))
	if err != nil {
		fmt.Println("配置文件获取内容失败", err)
		return
	}
	file_yaml := authConfig{}
	err1 := yaml.Unmarshal(content, &file_yaml)
	if err1 != nil {
		fmt.Println("解析配置文件失败", err1)
	}
	fmt.Printf("%v\n", file_yaml)
	//获取runtimeconfig
	info := apiclient.GetInitInfo_v2(file_yaml.Init, file_yaml.ApiServerIP, file_yaml.Token, file_yaml.CaHash, file_yaml.ConfigFileName)
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	cmd.Execute(runtimeConfig)
}
