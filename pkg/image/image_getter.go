package image

import (
	"fmt"
	"log"
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
		image = strings.Replace(image, k, v, -1)
	}
	pullCmd := fmt.Sprintf("docker pull %s", image)
	log.Printf("Using image %s\n", pullCmd)
	cmd := exec.Command("docker", "pull", image)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = tagImage(image, oldImage)
	if err != nil {
		return err
	}
	err = deleteImage(image)
	if err != nil {
		return err
	}
	return nil
}

func tagImage(oldImage string, newImage string) error {
	tagCmd := fmt.Sprintf("docker tag %s %s", oldImage, newImage)
	log.Printf("Tagging image: %s\n", tagCmd)
	cmd := exec.Command("docker", "tag", oldImage, newImage)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func deleteImage(image string) error {
	deleteCmd := fmt.Sprintf("docker rmi %s", image)
	log.Printf("Deleting Image: %s \v", deleteCmd)
	cmd := exec.Command("docker", "rmi", image)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
