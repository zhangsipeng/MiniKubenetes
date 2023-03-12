package iptables

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/data/WorkloadResources"
	"strconv"
)

const kubeServiceChain = "KUBE-SERVICE"

func (ipt *IPTable) Init() error {
	var err error
	err = ipt.clear()
	if err != nil {
		return err
	}
	err = ipt.AppendRule("nat", "POSTROUTING", "-o", "minik8s-vxlan0", "-j", "MASQUERADE")
	if err != nil {
		return err
	}
	err = ipt.AddChain("nat", kubeServiceChain)
	if err != nil {
		return nil
	}
	err = ipt.addReference("nat", "PREROUTING", kubeServiceChain)
	if err != nil {
		return err
	}
	err = ipt.addReference("nat", "OUTPUT", kubeServiceChain)
	return err
}

func (ipt *IPTable) AddService(service ServiceResources.Service, pods []WorkloadResources.Pod) error {
	podNum := len(pods)
	for i, port := range service.Spec.Ports {
		sha1Hash := sha1.Sum([]byte(service.Metadata.Name + "-" + strconv.Itoa(i)))
		finalName := base64.RawURLEncoding.EncodeToString(sha1Hash[:])[:15]
		chain := "K8S-SERVICE-" + finalName
		err := ipt.AddChain("nat", chain)
		if err != nil {
			return err
		}
		if port.NodePort != 0 {
			err = ipt.addReference("nat", kubeServiceChain, chain,
				"-p", "tcp", "--dport", strconv.Itoa(int(port.NodePort)))
			if err != nil {
				return err
			}
		}
		err = ipt.addReference("nat", kubeServiceChain, chain,
			"-p", "tcp", "-d", service.Spec.ClusterIP,
			"--dport", strconv.Itoa(int(port.Port)))
		if err != nil {
			return err
		}
		for j, pod := range pods {
			err = ipt.AppendRule("nat", chain, "-p", "tcp",
				"-m", "statistic", "--mode", "nth", "--every", strconv.Itoa(podNum-j), "--packet", "0",
				"-j", "DNAT", "--to-destination", pod.Spec.IP+":"+strconv.Itoa(int(port.TargetPort)))
		}
	}
	return nil
}

func (ipt *IPTable) GetRule(chainName string) ([]string, error) {
	buffer, err := ipt.getChainContent("nat", chainName)
	if err != nil {
		return nil, err
	}
	content := buffer.Bytes()
	rules := make([]string, 0)
	chain := bytes.Split(content, []byte{'\n'})
	for _, rule := range chain[2:] {
		rules = append(rules, string(rule))
	}
	return rules, nil
}

func (ipt *IPTable) GetChain() ([]string, error) {
	buffer, err := ipt.getTableContent("nat")
	if err != nil {
		return nil, err
	}
	content := buffer.Bytes()
	chains := make([]string, 0)
	table := bytes.Split(content, []byte("Chain "))
	for _, chain := range table {
		chainName := string(bytes.Split(chain, []byte(" "))[0])
		chains = append(chains, chainName)
	}
	return chains, nil
}

func (ipt *IPTable) clear() error {
	ipt.DeleteRule("nat", "POSTROUTING", "-o", "minik8s-vxlan0", "-j", "MASQUERADE")
	rules, err := ipt.GetRule("PREROUTING")
	for i := len(rules) - 1; i >= 0; i-- {
		if len(rules[i]) >= 4 && rules[i][0:4] == "KUBE" {
			err = ipt.deleteIndexRule("nat", "PREROUTING", i+1)
			if err != nil {
				return err
			}
		}
	}
	rules, err = ipt.GetRule("OUTPUT")
	for i := len(rules) - 1; i >= 0; i-- {
		if len(rules[i]) >= 4 && rules[i][0:4] == "KUBE" {
			err = ipt.deleteIndexRule("nat", "OUTPUT", i+1)
			if err != nil {
				return err
			}
		}
	}
	chains, err := ipt.GetChain()
	if err != nil {
		return err
	}
	for _, chain := range chains {
		if len(chain) >= 4 && chain[0:4] == "KUBE" {
			err := ipt.DeleteChain("nat", chain)
			if err != nil {
				return err
			}
		}
	}
	for _, chain := range chains {
		if len(chain) >= 3 && chain[0:3] == "K8S" {
			err := ipt.DeleteChain("nat", chain)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
