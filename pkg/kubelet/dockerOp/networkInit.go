package dockerOp

import (
	"context"
	"errors"
	"example/Minik8s/pkg/data/ClusterResources"
	"example/Minik8s/pkg/kubelet/dockerConst"
	"fmt"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func RemoveDockerNetworkIfExist(cli *client.Client) error {
	networks, err := cli.NetworkList(context.TODO(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: dockerConst.NetworkName,
		}),
	})
	if err != nil {
		return err
	}
	for _, network := range networks {
		err := cli.NetworkRemove(context.TODO(), network.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDockerNetwork(cli *client.Client, node ClusterResources.Node) error {
	podCIDR := node.Spec.PodCIDR
	podCIDRIP, podCIDRNet, err := net.ParseCIDR(podCIDR)
	if err != nil {
		return err
	}
	podGatewayIP := podCIDRIP.To4()
	if podGatewayIP == nil {
		return errors.New(fmt.Sprintf("%s: invalid IPv4 address", podCIDR))
	}
	if ones, _ := podCIDRNet.Mask.Size(); ones != 24 {
		return errors.New(fmt.Sprintf("%s: unrecognized pod CIDR network", podCIDR))
	}
	podGatewayIP[3] = 254
	_, err = cli.NetworkCreate(context.TODO(), dockerConst.NetworkName, types.NetworkCreate{
		CheckDuplicate: true,
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet:  podCIDR,
					Gateway: podGatewayIP.String(),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
