package vxlan

import (
	"errors"
	"net"
	"os/exec"
)

type peerInfo struct {
	remoteVxlanIPCidr, remotePodIPSubnet string
}

var addedPeer = map[string]peerInfo{}

func AddPeer(remoteIP, remoteVxlanIPCidr, remotePodIPSubnet string) error {
	if info, ok := addedPeer[remoteIP]; ok {
		if info.remoteVxlanIPCidr == remoteVxlanIPCidr &&
			info.remotePodIPSubnet == remotePodIPSubnet {
			return nil
		} else {
			return errors.New("changing peer info is not supported")
		}
	}
	if err := exec.Command("bridge", "fdb",
		"append", "00:00:00:00:00:00",
		"dev", VxlanIface,
		"dst", remoteIP).Run(); err != nil {
		return err
	}
	remoteVxlanIP, _, err := net.ParseCIDR(remoteVxlanIPCidr)
	if err != nil {
		return err
	}
	if err := exec.Command("ip",
		"route", "add", remotePodIPSubnet,
		"via", remoteVxlanIP.String()).Run(); err != nil {
		return err
	}
	addedPeer[remoteIP] = peerInfo{
		remoteVxlanIPCidr: remoteVxlanIPCidr,
		remotePodIPSubnet: remotePodIPSubnet,
	}
	return nil
}

func DelPeer(remoteIP, remoteVxlanIP, remotePodIPSubnet string) error {
	if err := exec.Command("ip",
		"route", "del", remotePodIPSubnet,
		"via", remoteVxlanIP).Run(); err != nil {
		return err
	}
	if err := exec.Command("bridge", "fdb",
		"del", "00:00:00:00:00:00",
		"dev", VxlanIface,
		"dst", remoteIP).Run(); err != nil {
		return err
	}
	return nil
}
