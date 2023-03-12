package iptables

import (
	"bytes"
	"os/exec"
	"strconv"
)

func New() (*IPTable, error) {
	path, err := exec.LookPath("iptables")
	if err != nil {
		return nil, err
	}
	ipt := &IPTable{path: path}
	return ipt, nil
}

func (ipt *IPTable) DeleteRule(table, chain string, ruleSpec ...string) error {
	args := append([]string{"-t", table, "-D", chain}, ruleSpec...)
	return ipt.run(args...)
}

func (ipt *IPTable) AddChain(table, chain string) error {
	return ipt.run("-t", table, "-N", chain)
}

func (ipt *IPTable) AppendRule(table, chain string, ruleSpec ...string) error {
	args := append([]string{"-t", table, "-A", chain}, ruleSpec...)
	return ipt.run(args...)
}

func (ipt *IPTable) ClearTable(table string) error {
	err := ipt.run("-t", table, "-F")
	if err != nil {
		return nil
	}
	err = ipt.run("-t", table, "-X")
	return err
}

func (ipt *IPTable) DeleteChain(table, chain string) error {
	var err error
	err = ipt.run("-t", table, "-F", chain)
	if err != nil {
		return err
	}
	err = ipt.run("-t", table, "-X", chain)
	return err
}

func (ipt *IPTable) getChainContent(table, chain string) (bytes.Buffer, error) {
	buffer, err := ipt.runWithOutPut("-t", table, "-L", chain)
	return buffer, err
}

func (ipt *IPTable) addReference(table string, srcChain, dstChain string, ruleSpec ...string) error {
	args := append([]string{"-t", table, "-A", srcChain, "-j", dstChain}, ruleSpec...)
	return ipt.run(args...)
}

func (ipt *IPTable) getTableContent(table string) (bytes.Buffer, error) {
	buffer, err := ipt.runWithOutPut("-t", table, "-L")
	return buffer, err
}

func (ipt *IPTable) deleteIndexRule(table, chain string, index int) error {
	return ipt.run("-t", table, "-D", chain, strconv.Itoa(index))
}

func (ipt *IPTable) run(args ...string) error {
	args = append([]string{ipt.path}, args...)
	cmd := exec.Cmd{
		Path: ipt.path,
		Args: args,
	}
	err := cmd.Run()
	return err
}

func (ipt *IPTable) runWithOutPut(args ...string) (bytes.Buffer, error) {
	args = append([]string{ipt.path}, args...)
	var buffer bytes.Buffer
	cmd := exec.Cmd{
		Path:   ipt.path,
		Args:   args,
		Stdout: &buffer,
	}
	err := cmd.Run()
	return buffer, err
}
