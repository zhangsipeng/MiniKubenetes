package dnsOp

import (
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/data/WorkloadResources"
	"fmt"
	"os"
	"os/exec"
)

var dnsServer string
var ddnsCommandArgs []string

func AddARecord(name, nameAddr string) (err error) {
	file, err := os.CreateTemp("", "*")
	if err != nil {
		return
	}
	_, err = file.WriteString(fmt.Sprintf(`server %s
update del %s IN A
update add %s 6000 IN A %s
send
quit
`, dnsServer, name, name, nameAddr))
	if err != nil {
		return
	}
	err = exec.Command("nsupdate", append(ddnsCommandArgs, file.Name())...).Run()
	os.Remove(file.Name())
	return nil
}

func AddPod(pod WorkloadResources.Pod) error {
	return AddARecord(
		fmt.Sprintf("%s.pod.minik8s", pod.Metadata.Name),
		pod.Spec.IP)
}

func AddService(service ServiceResources.Service) error {
	return AddARecord(
		fmt.Sprintf("%s.service.minik8s", service.Metadata.Name),
		service.Spec.ClusterIP)
}

func InitDnsOp(dnsServerIP string, ddnsKeyFilename *string) {
	dnsServer = dnsServerIP
	if ddnsKeyFilename == nil {
		ddnsCommandArgs = []string{}
	} else {
		ddnsCommandArgs = []string{"-k", *ddnsKeyFilename}
	}
}
