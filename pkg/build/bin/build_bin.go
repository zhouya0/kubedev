package bin

import (
	"fmt"
	"kubedev/pkg/env"
	"log"
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

func (i *BinConfig) SetEnv(cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_GIT_VERSION_FILE=%s", i.KubeGitVersionFile))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PLATFORMS=%s", i.KubeBuildPlatforms))
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
	b.KubeBuildPlatforms = k.BuildPlatform
}

func BuildBinary() error {
	// Step 1: init configuration file
	binConfig := NewDefaultBinConfig()
	mergeKubeDevConfigAndBinConfig(&env.Config, binConfig)

	// Step 2: generate version file
	env.WriteVersionFile(env.KubeVersionFile)

	// step 3: build binary
	cmd := exec.Command("bash", "build/run.sh", "make")
	cmd.Env = os.Environ()
	binConfig.SetEnv(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error building binary: %v", err)
		return err
	}
	log.Printf("Build binary: %s", string(out))
	return nil
}
