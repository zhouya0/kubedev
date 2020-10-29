package image

import (
	kubedevlog "kubedev/pkg/log"
	"log"
	"os/exec"
	"strings"
)

// Images is map for image convert special for CN users
var Images = map[string]string{
	"k8s.gcr.io": "daocloud.io/daocloud",
	"quay.io":    "quay.azk8s.cn",
	"us.gcr.io/k8s-artifacts-prod/build-image": "daocloud.io/daocloud",
	"k8s.gcr.io/build-image": "daocloud.io/daocloud",
}

// PullImage will try to pull image using CN source
func PullImage(image string, logger *log.Logger) error {
	oldImage := image
	if image == "" {
		return nil
	}
	splitTags := strings.Split(oldImage, ":")
	version := splitTags[1]
	repoTags := strings.Split(splitTags[0], "/")
	imageName := repoTags[len(repoTags)-1]
	repoTag := strings.Join(repoTags[:len(repoTags)-1], "/")
	newRepoTag := Images[repoTag]
	image = newRepoTag + "/" + imageName + ":" + version
	cmd := exec.Command("docker", "pull", image)
	out, err := cmd.CombinedOutput()
	logger.Printf("Pull image output: %s", string(out))
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	err = tagImage(image, oldImage)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	err = deleteImage(image)
	if err != nil {
		kubedevlog.LogErrorMessage(logger, err)
		return err
	}
	return nil
}

func tagImage(oldImage string, newImage string) error {
	cmd := exec.Command("docker", "tag", oldImage, newImage)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func deleteImage(image string) error {
	cmd := exec.Command("docker", "rmi", image)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
