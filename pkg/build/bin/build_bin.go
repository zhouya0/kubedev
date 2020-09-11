package bin

import (
	"fmt"
	"kubedev/pkg/cli"
	"kubedev/pkg/env"
	imageGetter "kubedev/pkg/image"
	kubedevlog "kubedev/pkg/log"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

var HardCodeString []string = []string{"KUBE_GIT_VERSION_FILE", "KUBE_BUILD_PLATFORMS", "GOFLAGS"}

type BinConfig struct {
	KubeGitVersionFile string
	KubeBuildPlatforms string
	GoFlags            string
	OverrideVersion    string
}

func (i *BinConfig) String() string {
	v := reflect.ValueOf(*i)
	count := v.NumField()
	binConfigString := ""
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if f.String() != "" {
			binConfigString = binConfigString + " " + HardCodeString[i] + "=" + f.String()
		}

	}
	return binConfigString
}

func (i *BinConfig) SetEnv(cmd *exec.Cmd, arch string) {
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_GIT_VERSION_FILE=%s", i.KubeGitVersionFile))
	if arch != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PLATFORMS=%s", arch))
	} else {
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PLATFORMS=%s", i.KubeBuildPlatforms))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOFLAGS=%s", i.GoFlags))
}

func NewDefaultBinConfig() *BinConfig {
	return &BinConfig{
		KubeGitVersionFile: env.KubeVersionFile,
		KubeBuildPlatforms: "linux/amd64",
		GoFlags:            "-tags=nokmem",
	}
}

func mergeKubeDevConfigAndBinConfig(k *env.KubeDevConfig, b *BinConfig) {
	if k.BuildPlatform != "" {
		b.KubeBuildPlatforms = k.BuildPlatform
	}
}

func BuildBinary(args []string, arch string) error {
	logger := kubedevlog.NewLogger()
	status := cli.NewStatus()

	// Step 1: init configuration file
	binConfig := NewDefaultBinConfig()
	mergeKubeDevConfigAndBinConfig(&env.Config, binConfig)

	// Step 2: generate version file
	status.Start(fmt.Sprintf("Writing version file %s", env.WriteIcon))
	err := env.WriteVersionFile(env.KubeVersionFile, env.Config.OverrideKubeVersion)
	status.End(err == nil)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step 3: pull all images
	status.Start(fmt.Sprintf("Pulling building images %s", env.ImageIcon))
	err = prePullImages(logger)
	status.End(err == nil)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}


	// step 4: build binary
	status.Start(fmt.Sprintf("Building binary %s %s", args[0], env.BuildIcon))
	cmd := exec.Command("bash", "build/run.sh", "make", args[0])
	cmd.Env = os.Environ()
	binConfig.SetEnv(cmd, arch)
	out, err := cmd.CombinedOutput()
	status.End(err == nil)
	logger.Printf("Build binary: %s", string(out))
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	fmt.Printf("Building binary %s success! File can be found in:\n %s\n", args[0], filepath.Join(env.KubeBinPath, arch, args[0]))
	return nil
}

func prePullImages(logger *log.Logger) error {
	kubeImages := env.GetAllImages()
	err := imageGetter.PullImage(kubeImages.KubeCross, logger)
	if err != nil {
		return err
	}
	return nil
}