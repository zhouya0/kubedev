package image

import (
	"fmt"
	"kubedev/pkg/cli"
	"kubedev/pkg/env"
	imageGetter "kubedev/pkg/image"
	kubedevlog "kubedev/pkg/log"
	"os"
	"os/exec"
	"reflect"
)

var HardCodeString []string = []string{"KUBE_FASTBUILD", "KUBE_BUILD_HYPERKUBE", "KUBE_BUILD_CONFORMANCE", "KUBE_DOCKER_IMAGE_TAG", "KUBE_DOCKER_REGISTRY", "KUBE_GIT_VERSION_FILE", "KUBE_BUILD_PULL_LATEST_IMAGES", "GOFLAGS"}

type ImageConfig struct {
	KubeFastBuild             string
	KubeBuildHyperkube        string
	KubeBuildConformance      string
	KubeDockerImageTag        string
	KubeDockerRegistry        string
	KubeGitVersionFile        string
	KubeBuildPullLatestImages string
	GoFlags                   string
}

func NewDefaultImageConfig() *ImageConfig {
	return &ImageConfig{
		KubeFastBuild:             "true",
		KubeBuildHyperkube:        "n",
		KubeBuildConformance:      "n",
		KubeDockerImageTag:        "",
		KubeDockerRegistry:        "",
		KubeGitVersionFile:        env.KubeVersionFile,
		KubeBuildPullLatestImages: "n",
		GoFlags:                   "-tags=nokmem",
	}
}

func (i *ImageConfig) String() string {
	v := reflect.ValueOf(*i)
	count := v.NumField()
	imageConfigString := ""
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if f.String() != "" {
			imageConfigString = imageConfigString + " " + HardCodeString[i] + "=" + f.String()
		}

	}
	return imageConfigString
}

func (i *ImageConfig) SetEnv(cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_FASTBUILD=%s", i.KubeFastBuild))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_HYPERKUBE=%s", i.KubeBuildHyperkube))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_CONFORMANCE=%s", i.KubeBuildConformance))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_DOCKER_IMAGE_TAG=%s", i.KubeDockerImageTag))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_DOCKER_REGISTRY=%s", i.KubeDockerRegistry))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_GIT_VERSION_FILE=%s", i.KubeGitVersionFile))
	cmd.Env = append(cmd.Env, fmt.Sprintf("KUBE_BUILD_PULL_LATEST_IMAGES=%s", i.KubeBuildPullLatestImages))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOFLAGS=%s", i.GoFlags))
}

func (i *ImageConfig) SetKubeDockerImageTag(s string) {
	i.KubeDockerImageTag = s
}

func (i *ImageConfig) SetKubeDockerRegistry(s string) {
	i.KubeDockerRegistry = s
}

func (i *ImageConfig) SetKubeGitVersionFile(s string) {
	i.KubeGitVersionFile = s
}

func mergeKubeDevConfigAndImageConfig(k *env.KubeDevConfig, i *ImageConfig) {
	i.KubeDockerImageTag = k.DockerTag
	i.KubeDockerRegistry = k.DockerRegistry
}

func BuildImages(args []string) error {
	logger := kubedevlog.NewLogger()
	status := cli.NewStatus()

	// Step 1: init config file
	imageConfig := NewDefaultImageConfig()
	mergeKubeDevConfigAndImageConfig(&env.Config, imageConfig)

	// Step 2: pull all images
	status.Start(fmt.Sprintf("Pulling building images %s", env.ImageIcon))
	err := prePullImages()
	status.End(err == nil)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step 3: generate version file
	status.Start(fmt.Sprintf("Writing version file %s", env.WriteIcon))
	err = env.WriteVersionFile(env.KubeVersionFile, env.Config.OverrideKubeVersion)
	status.End(err == nil)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	// Step 4: make release
	status.Start(fmt.Sprintf("Making release images %s", env.BuildIcon))
	cmd := exec.Command("make", "release-images")
	cmd.Env = os.Environ()
	imageConfig.SetEnv(cmd)
	out, err := cmd.CombinedOutput()
	status.End(err == nil)
	logger.Printf("Build images: %s", string(out))
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}

	fmt.Printf("Building images success!")
	return nil
}

func prePullImages() error {
	kubeImages := env.GetAllImages()
	err := imageGetter.PullImage(kubeImages.DebianBase)
	if err != nil {
		return err
	}
	// imageGetter.PullImage(kubeImages.DebianHyperKubeBase)
	err = imageGetter.PullImage(kubeImages.KubeCross)
	if err != nil {
		return err
	}
	err = imageGetter.PullImage(kubeImages.KubePause)
	if err != nil {
		return err
	}
	err = imageGetter.PullImage(kubeImages.DebianIptables)
	if err != nil {
		return err
	}
	return nil
}
