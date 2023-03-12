package vxlan

import (
	"fmt"
	"os/exec"
)

const vxlanId = 42

func InitVxlan(localVxlanIPCidr string) error {
	if err := exec.Command("ip", "link",
		"add", VxlanIface, "type", "vxlan",
		"id", fmt.Sprint(vxlanId), "dstport", "0").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "addr",
		"add", localVxlanIPCidr, "dev", VxlanIface).Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link",
		"set", "up", "dev", VxlanIface).Run(); err != nil {
		return err
	}
	return nil
}

func DelVxlanIfExists() error {
	_ = exec.Command("ip", "link", "delete", VxlanIface).Run()
	return nil
}
