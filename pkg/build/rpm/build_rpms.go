package rpm

import (
	"fmt"
	"kubedev/pkg/build/bin"
	"kubedev/pkg/build/rpm/files"
	"kubedev/pkg/cli"
	"kubedev/pkg/env"
	kubedevlog "kubedev/pkg/log"
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
	rpmSrpms  string = "rpmbuild/SRPMS"
	rpmBuild  string = "rpmbuild/BUILD"
)

func BuildRPM(args []string, arch string) error {
	kubeversion := env.NewKubeVerisonOverride(env.Config.OverrideKubeVersion)
	logger := kubedevlog.NewLogger()
	status := cli.NewStatus()

	// Step1: Build kubelet binary file
	currentDir, _ := os.Getwd()
	componentDir := filepath.Join(currentDir, env.KubeBinPath, arch, args[0])
	if !util.CheckExist(componentDir) {
		err := bin.BuildBinary(args, arch)
		if err != nil {
			return err
		}
	}

	status.Start(fmt.Sprintf("Packaging binary to RPM %s", env.PackageIcon))
	// Step2: Make directories for RPM build use
	makeRPMBuildDir(logger)

	// Step2: Move kubelet to Source directory
	err := cpComponentToSource(arch, args[0], logger, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step3: Write pre rpmbuild files
	err = writePreBuildFiles(args[0], logger, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step4: tar the source file
	err = tarSourceFile(args[0], logger, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step5: rpm build
	var target string
	if arch == "linux/arm64" {
		target = "aarch64"
	}
	err = RPMBuild(args[0], logger, kubeversion, target)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	status.End(err == nil)

	fmt.Printf("Building RPM %s success! Package can be found in: \n %s\n", args[0], filepath.Join(util.GetHomeDir(), rpmRpms))
	return nil
}

func RPMBuild(component string, logger *log.Logger, kubeversion env.KubeVersion, target string) error {
	specFile := filepath.Join(util.GetHomeDir(), rpmSpecs, component+".spec")
	versionDefine := fmt.Sprintf("_version %s", env.GetKubeVersionNoV(kubeversion))
	// TODO: what release should be used here?
	releaseDefine := fmt.Sprintf("_release %s", "00")
	cmd := &exec.Cmd{}
	if target != "" {
		cmd = exec.Command("rpmbuild", "-ba", specFile, "--define", versionDefine, "--define", releaseDefine, "--target", target)
	} else {
		cmd = exec.Command("rpmbuild", "-ba", specFile, "--define", versionDefine, "--define", releaseDefine)
	}

	out, err := cmd.CombinedOutput()
	logger.Println(string(out))
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	return nil
}

func makeRPMBuildDir(logger *log.Logger) {
	pathRPM := filepath.Join(util.GetHomeDir(), rpmRpms)
	pathSRPM := filepath.Join(util.GetHomeDir(), rpmSrpms)
	pathSPEC := filepath.Join(util.GetHomeDir(), rpmSpecs)
	pathSOURCE := filepath.Join(util.GetHomeDir(), rpmSource)
	pathBUILD := filepath.Join(util.GetHomeDir(), rpmBuild)
	if !util.CheckExist(pathRPM) {
		os.MkdirAll(pathRPM, os.ModePerm)
	}
	if !util.CheckExist(pathSRPM) {
		os.MkdirAll(pathSRPM, os.ModePerm)
	}
	if !util.CheckExist(pathSPEC) {
		os.MkdirAll(pathSPEC, os.ModePerm)
	}
	if !util.CheckExist(pathSOURCE) {
		os.MkdirAll(pathSOURCE, os.ModePerm)
	}
	if !util.CheckExist(pathBUILD) {
		os.MkdirAll(pathBUILD, os.ModePerm)
	}
}

func cpComponentToSource(arch string, component string, logger *log.Logger, kubeversion env.KubeVersion) error {
	componentDir := filepath.Join(env.KubeBinPath, arch, component)
	rpmSourceDir := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component)
	// cmd := exec.Command("cp", "-p", componentDir, rpmSourceDir)
	err := util.CopyFile(componentDir, rpmSourceDir)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	// change the mod of kubelet to 0755
	os.Chmod(rpmSourceDir, 0755)

	return nil
}

func writePreBuildFiles(component string, logger *log.Logger, kubeversion env.KubeVersion) error {
	err := writeComponentSpec(component, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	err = writeComponentService(component, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	err = writeComponentEnv(component, kubeversion)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	return nil
}

func writeComponentSpec(component string, kubeversion env.KubeVersion) error {
	rpmSourceSpecFile := filepath.Join(util.GetHomeDir(), rpmSpecs, component+".spec")
	// case component == "kubelet":
	err := util.WriteFile(rpmSourceSpecFile, files.KubeletSpec)
	if err != nil {
		return err
	}
	return nil
}

func writeComponentService(component string, kubeversion env.KubeVersion) error {
	rpmServiceFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component+".service")
	err := util.WriteFile(rpmServiceFile, files.KubeletService)
	if err != nil {
		return err
	}
	return nil
}

func writeComponentEnv(component string, kubeversion env.KubeVersion) error {
	rpmEnvFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion), component+".env")
	err := util.WriteFile(rpmEnvFile, files.KubeletEnv)
	if err != nil {
		return err
	}
	return nil
}

func tarSourceFile(component string, logger *log.Logger, kubeversion env.KubeVersion) error {
	sourcePath := filepath.Join(util.GetHomeDir(), rpmSource)
	tarFile := filepath.Join(util.GetHomeDir(), rpmSource, env.GetComponentDirName(component, kubeversion)+".tar.gz")
	cmd := exec.Command("tar", "-czvf", tarFile, "-C", sourcePath, env.GetComponentDirName(component, kubeversion))
	out, err := cmd.CombinedOutput()
	logger.Println(string(out))
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	return nil
}
