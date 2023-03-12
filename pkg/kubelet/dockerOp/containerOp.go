package dockerOp

import (
	"context"
	"example/Minik8s/pkg/kubelet/dockerConst"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func ContainerRemove(cli *client.Client, podFilter filters.Args) error {
	containerList, err := cli.ContainerList(context.TODO(),
		types.ContainerListOptions{All: true, Filters: podFilter})
	if err != nil {
		return err
	}
	volumeList, err := cli.VolumeList(context.TODO(), podFilter)
	if err != nil {
		return err
	}
	for _, containerItem := range containerList {
		err := cli.ContainerStop(context.TODO(), containerItem.ID, nil)
		if err != nil {
			return err
		}
	}
	for _, containerItem := range containerList {
		err := cli.ContainerRemove(context.TODO(), containerItem.ID,
			types.ContainerRemoveOptions{
				RemoveVolumes: true,
				RemoveLinks:   false,
				Force:         false,
			})
		if err != nil {
			return err
		}
	}
	for _, volumeItem := range volumeList.Volumes {
		err := cli.VolumeRemove(context.TODO(), volumeItem.Name, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func ContainerRemoveAllK8s(cli *client.Client) error {
	podFilter := filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: dockerConst.LabelMinik8s,
		},
	)
	return ContainerRemove(cli, podFilter)
}
