package image

import (
	"fmt"
	"os/exec"
	"strings"
)

// Images is map for image convert special for CN users
var Images = map[string]string{
	"gcr.io":     "gcr.azk8s.cn",
	"k8s.gcr.io": "gcr.azk8s.cn/google-containers",
	"quay.io":    "quay.azk8s.cn",
}

// PullImage will try to pull image using CN source
func PullImage(image string) error {
	oldImage := image
	for k, v := range Images {
		strings.Replace(image, k, v, -1)
	}
	pullCmd := fmt.Sprintf("docker pull %s", image)
	cmd := exec.Command(pullCmd)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = tagImage(image, oldImage)
	if err != nil {
		return err
	}
	err = deleteImage(oldImage)
	if err != nil {
		return err
	}
	return nil
}

func tagImage(oldImage string, newImage string) error {
	tagCmd := fmt.Sprintf("docker tag %s %s", oldImage, newImage)
	cmd := exec.Command(tagCmd)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func deleteImage(image string) error {
	deleteCmd := fmt.Sprintf("docker rmi %s", image)
	cmd := exec.Command(deleteCmd)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
