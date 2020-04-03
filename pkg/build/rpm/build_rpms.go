package rpm

import (
	"fmt"
	"kubedev/pkg/build/bin"
	"kubedev/pkg/env"
	"kubedev/pkg/util"
	"log"
	"os/exec"
)

func BuildRPM(args []string, arch string) error {
	// Step1: Build kubelet binary file
	if !util.CheckfileExist(env.KubeBinPath + "/" + arch + "/" + args[0]) {
		err := bin.BuildBinary(args, arch)
		if err != nil {
			return err
		}
		fmt.Println("Magic. It doesn't exist!", env.KubeBinPath+"/"+arch+"/"+args[0])
	}

	// Step2: Move kubelet to Source directory
	err := cpKubeletToSource(arch, args[0])
	if err != nil {
		return err
	}

	// Step3: init kubelet service and env file

	return nil
}

func cpKubeletToSource(arch string, component string) error {
	cmd := exec.Command("cp", "-p", env.KubeBinPath+"/"+arch+"/"+component, util.GetHomeDir()+"/"+"SOURCES/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Printf("Copy component %s succeed: %s", component, string(out))
	return nil
}
