package bin

import (
	"fmt"
	"kubedev/pkg/cli"
	"kubedev/pkg/env"
	kubedevlog "kubedev/pkg/log"
	"os"
	"os/exec"
	"reflect"
)

var HardCodeString []string = []string{"KUBE_GIT_VERSION_FILE", "KUBE_BUILD_PLATFORMS", "GOFLAGS"}

type BinConfig struct {
	KubeGitVersionFile string
	KubeBuildPlatforms string
	GoFlags            string
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
	icon := "📝"
	status.Start(fmt.Sprintf("Writing version file %s", icon))
	err := env.WriteVersionFile(env.KubeVersionFile)
	status.End(err == nil)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// step 3: build binary
	icon = "🔨"
	status.Start(fmt.Sprintf("Building binary %s %s", args[0], icon))
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
	return nil
}
