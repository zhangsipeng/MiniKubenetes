package docker

import (
	"context"
	"example/Minik8s/pkg/const/urlconst"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/serverless/data"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func createBuildContext(serverlessService *Serverless.Service) (tarfile *os.File, err error) {
	prevdir, err := os.Getwd()
	if err != nil {
		return
	}
	gitDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return
	}
	defer func() {
		os.Chdir(prevdir)
		os.RemoveAll(gitDir)
	}()
	err = exec.Command("git", "clone", "--depth=1", serverlessService.Spec.GitUrl, gitDir).Run()
	if err != nil {
		return
	}
	os.Chdir(gitDir)
	tarfile, err = os.CreateTemp("", "*")
	if err != nil {
		return
	}
	tarfileName := tarfile.Name()
	tarfile.Close()
	err = exec.Command("git", "archive", "--format=tar", "-o", tarfileName, "HEAD").Run()
	if err != nil {
		return
	}
	tarfile, err = os.Open(tarfileName)
	return
}

func GenerateDockerImage(serverlessService Serverless.Service) (imageName string, err error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return
	}
	tarfile, err := createBuildContext(&serverlessService)
	if err != nil {
		return
	}
	defer func() {
		tarfile.Close()
		os.Remove(tarfile.Name())
	}()
	tagName := fmt.Sprintf("%s:%d/%s",
		data.GetCredential().YamlConfig.APIServerIP,
		urlconst.PortDockerRegistry,
		serverlessService.Metadata.Name)
	imageBuildResponse, err := cli.ImageBuild(context.TODO(), tarfile, types.ImageBuildOptions{
		Tags: []string{tagName},
	})
	if err != nil {
		return
	}
	defer imageBuildResponse.Body.Close()
	if err != nil {
		return
	}
	resody, _ := ioutil.ReadAll(imageBuildResponse.Body)
	log.Println(string(resody))
	cli.ImagePush(context.TODO(), tagName, types.ImagePushOptions{})
	return tagName, nil
}
