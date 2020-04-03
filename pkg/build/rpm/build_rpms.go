package rpm

import (
	"fmt"
	"kubedev/pkg/build/bin"
	"kubedev/pkg/build/rpm/files"
	"kubedev/pkg/env"
	"kubedev/pkg/util"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	rpmSource string = "rpmbuild/SOURCES"
	rpmSpecs  string = "rpmbuild/SPECS"
	rpmRpms   string = "rpmbuild/RPMS"
)

var kubeversion env.KubeVersion = env.NewKubeVersion()

func BuildRPM(args []string, arch string) error {
	// Step1: Build kubelet binary file
	currentDir, _ := os.Getwd()
	componentDir := filepath.Join(currentDir, env.KubeBinPath, arch, args[0])
	if !util.CheckExist(componentDir) {
		err := bin.BuildBinary(args, arch)
		if err != nil {
			return err
		}
		fmt.Println("Magic. It doesn't exist!", componentDir)
	}
	// Step2: Move kubelet to Source directory
	err := cpComponentToSource(arch, args[0])
	if err != nil {
		return err
	}

	// Step3: Write pre rpmbuild files
	err = writePreBuildFiles(args[0])
	if err != nil {
		return err
	}

	// Step4: tar the source file
	err = tarSourceFile(args[0])
	if err != nil {
		return err
	}

	// Step5: rpm build
	err = RPMBuild(args[0])
	if err != nil {
		return err
	}

	return nil
}

func RPMBuild(component string) error {
	specFile := filepath.Join(util.GetHomeDir(), rpmSpecs, component+".spec")
	versionDefine := fmt.Sprintf("_version %s", env.GetKubeVersionNoV(kubeversion))
	// TODO: what release should be used here?
	releaseDefine := fmt.Sprintf("_release %s", "00")
	cmd := exec.Command("rpmbuild", "-ba", specFile, "--define", versionDefine, "--define", releaseDefine)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error when rpmbuild: %s\n", string(out))
		return err
	}
	log.Printf("Building rpm success!")
	return nil
}

func cpComponentToSource(arch string, component string) error {
	componentDir := filepath.Join(env.KubeBinPath, arch, component)
	rpmSourceDir := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component)
	// cmd := exec.Command("cp", "-p", componentDir, rpmSourceDir)
	err := util.CopyFile(componentDir, rpmSourceDir)
	if err != nil {
		log.Printf("Error when copying: %v\n", err.Error())
		return err
	}
	// change the mod of kubelet to 0755
	os.Chmod(rpmSourceDir, 0755)

	return nil
}

func writePreBuildFiles(component string) error {
	err := writeComponentSpec(component)
	if err != nil {
		return err
	}

	err = writeComponentService(component)
	if err != nil {
		return err
	}

	err = writeComponentEnv(component)
	if err != nil {
		return err
	}
	return nil
}

func writeComponentSpec(component string) error {
	rpmSourceSpecFile := filepath.Join(util.GetHomeDir(), rpmSpecs, component+".spec")
	// case component == "kubelet":
	err := util.WriteFile(rpmSourceSpecFile, files.KubeletSpec)
	if err != nil {
		return err
	}
	return nil
}

func writeComponentService(component string) error {
	rpmServiceFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component+".service")
	err := util.WriteFile(rpmServiceFile, files.KubeletService)
	if err != nil {
		return err
	}
	return nil
}

func writeComponentEnv(component string) error {
	rpmEnvFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component+".env")
	err := util.WriteFile(rpmEnvFile, files.KubeletEnv)
	if err != nil {
		return err
	}
	return nil
}

func tarSourceFile(component string) error {
	sourcePath := filepath.Join(util.GetHomeDir(), rpmSource)
	tarFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion)+".tar.gz")
	cmd := exec.Command("tar", "-czvf", tarFile, "-C", sourcePath, env.GetComponentDirName(component, kubeversion))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error when tar source file: %s\n", string(out))
		return err
	}
	return nil
}
