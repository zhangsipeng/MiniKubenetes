package main

import "example/Minik8s/pkg/kubeproxy/iptables"

func main() {
	ipt, _ := iptables.New()
	//buffer := ipt.GetChain("nat")
	//content := buffer.Bytes()
	//table := bytes.Split(content, []byte("Chain "))
	//for _, chain := range table {
	//	name := bytes.Split(chain, []byte(" "))[0]
	//	str := string(name)
	//	println(str)
	//}
	err := ipt.Init(true)
	if err != nil {
		println(err)
	}
}
