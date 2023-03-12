package watchAPIServer

import (
	"context"
	"encoding/json"
	"errors"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/kubelet/dockerConst"
	"example/Minik8s/pkg/kubelet/dockerOp"
	"example/Minik8s/pkg/kubelet/reportStatus/reportPodStatus"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func deletePod(cli *client.Client, podName string, runtimeConfig runtimedata.RuntimeConfig) {
	log.Printf("pod %s: stopping...", podName)
	podFilter := filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: "pod=" + podName,
		},
	)
	err := dockerOp.ContainerRemove(cli, podFilter)
	if err != nil {
		panic(err)
	}
	reportPodStatus.ReportPodStatus(podName, "stopped", runtimeConfig)
	log.Printf("pod %s: stopped.", podName)
}

func putPod(cli *client.Client, podName string, runtimeConfig runtimedata.RuntimeConfig) {
	containerList, err := cli.ContainerList(context.TODO(),
		types.ContainerListOptions{
			Filters: filters.NewArgs(
				filters.KeyValuePair{
					Key:   "label",
					Value: dockerConst.LabelMinik8s + ",pod=" + podName,
				},
			),
		})
	if err != nil {
		panic(err)
	}
	if len(containerList) != 0 {
		panic("changing Pod is not supported yet!")
	}

	var pod WorkloadResources.Pod
	if err := json.Unmarshal(
		apiclient.Request(runtimeConfig,
			fmt.Sprintf("/api/v1/pods/%s", podName),
			nil, "GET"),
		&pod); err != nil {
		panic(err)
	}

	log.Printf("pod %s: starting...", podName)
	for _, containerItem := range pod.Spec.Container {
		imageName := containerItem.Image
		if imageList, err := cli.ImageList(context.TODO(),
			types.ImageListOptions{
				Filters: filters.NewArgs(
					filters.KeyValuePair{
						Key:   "reference",
						Value: imageName,
					},
				),
			}); err != nil {
			panic(err)
		} else if len(imageList) == 0 {
			io, err := cli.ImagePull(context.TODO(), imageName,
				types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}
			defer io.Close()
			iocontent, err := ioutil.ReadAll(io)
			if err != nil {
				panic(err)
			}
			log.Print(string(iocontent))
		}
	}

	labels := map[string]string{
		"pod":                       podName,
		"restartPoicy":              pod.Spec.RestartPolicy,
		dockerConst.LabelMinik8sKey: dockerConst.LabelMinik8sVal,
	}

	claimedVolume := map[string]string{}
	for _, volumeClaim := range pod.Spec.Volumes {
		volumeName := apiclient.MangleName("volume", podName, volumeClaim.Name)
		claimedVolume[volumeClaim.Name] = volumeName
		_, err := cli.VolumeCreate(context.TODO(),
			volume.VolumeCreateBody{
				Name:   volumeName,
				Labels: labels,
			})
		if err != nil {
			panic(err)
		}
	}
	for _, containerItem := range pod.Spec.Container {
		for _, m := range containerItem.VolumeMounts {
			if _, ok := claimedVolume[m.Name]; !ok {
				panic(errors.New(
					fmt.Sprintf("pod %s: container %s uses undeclared volume %s",
						podName, containerItem.Name, m.Name),
				))
			}
		}
	}
	restartPolicy := container.RestartPolicy{
		Name: pod.Spec.RestartPolicy,
	}
	podIP := pod.Spec.IP
	portBinding := make(nat.PortMap)
	for _, containerItem := range pod.Spec.Container {
		for _, portSpec := range containerItem.Ports {
			port, err := nat.NewPort(portSpec.Protocal, fmt.Sprint(portSpec.ContainerPort))
			if err != nil {
				panic(err)
			}
			var hostPort []nat.PortBinding
			if portSpec.HostPort != 0 {
				hostPort = append(hostPort, nat.PortBinding{
					HostPort: fmt.Sprint(portSpec.HostPort),
				})
			}
			if _, ok := portBinding[port]; ok {
				panic("conflict port map")
			}
			portBinding[port] = hostPort
		}
	}

	pauseContainer, err := cli.ContainerCreate(context.TODO(),
		&container.Config{
			Image:  dockerConst.PauseImage,
			Labels: labels,
		},
		&container.HostConfig{
			IpcMode:       container.IpcMode("shareable"),
			RestartPolicy: restartPolicy,
			PortBindings:  portBinding,
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				dockerConst.NetworkName: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: podIP,
					},
					IPAddress: podIP,
				},
			},
		}, nil, apiclient.MangleName("pause", podName),
	)
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(context.TODO(), pauseContainer.ID,
		types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	pauseContainerRef := "container:" + pauseContainer.ID
	for _, containerItem := range pod.Spec.Container {
		containerConfig := &container.Config{
			Image:  containerItem.Image,
			Labels: labels,
		}
		if len(containerItem.Command) != 0 {
			containerConfig.Cmd = containerItem.Command
		}

		mountInfo := []mount.Mount{}
		for _, m := range containerItem.VolumeMounts {
			mountInfo = append(mountInfo, mount.Mount{
				Source: claimedVolume[m.Name],
				Target: m.MountPath,
				Type:   mount.TypeVolume,
			})
		}

		createdContainer, err := cli.ContainerCreate(context.TODO(),
			containerConfig,
			&container.HostConfig{
				NetworkMode:   container.NetworkMode(pauseContainerRef),
				IpcMode:       container.IpcMode(pauseContainerRef),
				PidMode:       container.PidMode(pauseContainerRef),
				Mounts:        mountInfo,
				RestartPolicy: restartPolicy,
			},
			nil, nil,
			apiclient.MangleName("container", podName, containerItem.Name))
		if err != nil {
			panic(err)
		}
		err = cli.ContainerStart(context.Background(), createdContainer.ID,
			types.ContainerStartOptions{})
		if err != nil {
			panic(err)
		}
	}
	reportPodStatus.ReportPodStatus(podName, "running", runtimeConfig)
	log.Printf("pod %s: started.", podName)
}

func operateContainer(cli *client.Client, event watch.WatchEvent, runtimeConfig runtimedata.RuntimeConfig) {
	// TODO: check validity of podName
	eventKeyParts := strings.Split(event.Key, "/")
	podName := eventKeyParts[len(eventKeyParts)-1]
	log.Printf("pod %s: event %s", podName, event.Type)
	switch event.Type {
	case "DELETE":
		go deletePod(cli, podName, runtimeConfig)
		break
	case "PUT":
		go putPod(cli, podName, runtimeConfig)
		break
	}
}

func WatchPods(cli *client.Client, runtimeConfig runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	targetURL := fmt.Sprintf("/api/v1/podInNode/%s/pods/",
		runtimeConfig.YamlConfig.Others["nodeId"])
	log.Println(targetURL)
	go apiclient.WatchAPIWithRelativePath(runtimeConfig, targetURL, event)
	for {
		operateContainer(cli, <-event, runtimeConfig)
	}
}
