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
	slices := strings.Split(image, "/")
	slices[0] = Images[slices[0]]
	image = strings.Join(slices, "/")
	cmd := exec.Command("docker", "pull", image)
	out, err := cmd.CombinedOutput()
	log.Printf("Pull image output: %s", string(out))
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
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func deleteImage(image string) error {
	deleteCmd := fmt.Sprintf("docker rmi %s", image)
	log.Printf("Deleting Image: %s \v", deleteCmd)
	cmd := exec.Command("docker", "rmi", image)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
