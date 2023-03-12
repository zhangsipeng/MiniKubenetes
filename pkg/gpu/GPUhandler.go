package gpu

import (
	"bytes"
	"errors"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func sshClient(passwd, user, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
}

func sshExec(client *ssh.Client, remoteCmd string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b
	err = session.Run(remoteCmd)
	return b.Bytes(), err
}

func sshCopyToRemote(client *ssh.Client, srcFilePath, dstFilePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	srcFile, err := os.OpenFile(srcFilePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	session.Stdin = srcFile
	return session.Run(fmt.Sprintf(`cat > "%s"`, dstFilePath))
}

func sshCopyToLocal(client *ssh.Client, srcFilePath, dstFilePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	content, err := session.Output(fmt.Sprintf(`cat "%s"`, srcFilePath))
	return os.WriteFile(dstFilePath, content, 0600)
}

const (
	sshPasswd = "@@Bb^BW5"
	sshUser   = "stu634"
	sshHost   = "login.hpc.sjtu.edu.cn"
	scpHost   = "data.hpc.sjtu.edu.cn"
)

func RunJob(job WorkloadResources.GPUJob) []byte {
	prevdir, _ := os.Getwd()
	tempDir, _ := os.MkdirTemp("", "")
	defer func() {
		_ = os.Chdir(prevdir)
		_ = os.RemoveAll(tempDir)
	}()
	_ = os.Chdir(tempDir)
	jobName := job.Metadata.Name
	log.Printf("job %s in dir %s\n", jobName, tempDir)
	err := exec.Command("git", "clone", "--depth=1", job.Giturl, "git").Run()
	if err != nil {
		log.Fatalln("Your giturl is not working")
	}

	slurmContent := fmt.Sprintf(`#! /bin/bash

#SBATCH --job-name=%s
#SBATCH --partition=dgx2
#SBATCH --output=../stdout
#SBATCH --error=../stderr
#SBATCH -N %d
#SBATCH --ntasks-per-node=%d
#SBATCH --cpus-per-task=%d
#SBATCH --gres=gpu:%d

ulimit -s unlimited
ulimit -l ulimited
module load cuda/10.1.243-gcc-8.3.0

./entrypoint.sh
`,
		jobName, job.Config.Nodes, job.Config.Tskpernode,
		job.Config.Cpupertsk, job.Config.Gpunum)
	os.WriteFile("job.slurm", []byte(slurmContent), 0700)
	{
		os.Chdir("git")
		err = exec.Command("git", "archive", "--format=tar", "--prefix=git/",
			"-o", "../git.tar", "HEAD").Run()
		os.Chdir("..")
		if err != nil {
			panic(err)
		}
		defer os.Remove("git.tar")
		dataSshClient, err := sshClient(sshPasswd, sshUser, scpHost)
		if err != nil {
			panic(err)
		}
		defer dataSshClient.Close()
		_, err = sshExec(dataSshClient, fmt.Sprintf(`mkdir -p "%s"`, jobName))
		if err != nil {
			panic(err)
		}
		log.Println("ssh mkdir done")
		err = sshCopyToRemote(dataSshClient, "job.slurm", jobName+"/job.slurm")
		if err != nil {
			panic(err)
		}
		log.Println("scp job.slurm done")
		err = sshCopyToRemote(dataSshClient, "git.tar", jobName+"/git.tar")
		if err != nil {
			panic(err)
		}
		log.Println("scp git.tar done")
		_, err = sshExec(dataSshClient, fmt.Sprintf(`cd "%s"; tar -xf git.tar; rm git.tar`, jobName))
		if err != nil {
			panic(err)
		}
		log.Println("ssh untar git.tar done")
	}
	loginSshClient, err := sshClient(sshPasswd, sshUser, sshHost)
	if err != nil {
		panic(err)
	}
	defer loginSshClient.Close()
	response, err := sshExec(loginSshClient,
		fmt.Sprintf(`cd "%s/git"; sbatch ../job.slurm`, jobName))
	if err != nil {
		panic(errors.New(fmt.Sprintf("%s: %s", err.Error(), string(response))))
	}
	log.Println("submit job response: ", string(response))
	batchId := string(response)
	if idx := strings.Index(batchId, "Submitted batch job "); idx == -1 {
		panic(errors.New("cannot find batch id"))
	} else {
		batchId = strings.TrimSpace(batchId[idx+len("Submitted batch job "):])
	}
	log.Println("batch id is", batchId)
	for {
		out, err := sshExec(loginSshClient, fmt.Sprintf(`sacct | grep "^%s"`, batchId))
		if err != nil {
			panic(err)
		}
		Outstring := string(out)
		log.Printf(Outstring)
		if strings.Contains(Outstring, "COMPLETED") {
			dataSshClient, err := sshClient(sshPasswd, sshUser, scpHost)
			if err != nil {
				panic(err)
			}
			defer dataSshClient.Close()
			err = sshCopyToLocal(dataSshClient, jobName+"/stdout", "stdout")
			if err != nil {
				panic(err)
			}
			result, err := ioutil.ReadFile("stdout")
			if err != nil {
				panic(err)
			}
			return result
		}
		time.Sleep(10 * time.Second)
	}
}

func SubmitJob(job WorkloadResources.GPUJob, credential runtimedata.RuntimeConfig) {
	job.Res = string(RunJob(job))
	job.Phase = "done"
	_ = apiclient.Request(credential, "/api/v1/jobs", job, "PUT")
	return
}
