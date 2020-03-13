package image

import (
	"kubedev/pkg/env"
	imageGetter "kubedev/pkg/image"
	"os/exec"
)

// KUBE_FASTBUILD=false \
// KUBE_BUILD_HYPERKUBE=n \
// KUBE_BUILD_CONFORMANCE=n \
// KUBE_DOCKER_IMAGE_TAG=${KUBE_VERSION} \
// KUBE_DOCKER_REGISTRY=${HUB_PREFIX} \
// KUBE_GIT_VERSION_FILE=./dce_version \
// GOFLAGS="-tags=nokmem" \
// make release-images

type ImageConfig struct {
	KubeFastBuild        bool
	KubeBuildHyperkube   string
	KubeBuildConformance string
	KubeDockerImageTag   string
	KubeDockerRegistry   string
	KubeGitVersionFile   string
	GoFlags              string
}

func NewDefaultImageConfig() *ImageConfig {
	return &ImageConfig{
		KubeFastBuild:        false,
		KubeBuildHyperkube:   "n",
		KubeBuildConformance: "n",
		KubeDockerImageTag:   "",
		KubeDockerRegistry:   "",
		KubeGitVersionFile:   "",
	}
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

func BuildImages() error {
	cmd := exec.Command("make", "release-images")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func prePullImages() {
	kubeImages := env.GetAllImages()
	imageGetter.PullImage(kubeImages.DebianBase)
	imageGetter.PullImage(kubeImages.DebianHyperKubeBase)
	imageGetter.PullImage(kubeImages.KubeCross)
	imageGetter.PullImage(kubeImages.KubePause)
	imageGetter.PullImage(kubeImages.DebianIptables)
}
