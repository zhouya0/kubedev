package env

import (
	"log"
	"os/exec"
	"strings"
)

type KubeBuildImages struct {
	KubeCross           string
	KubePause           string
	DebianHyperKubeBase string
	DebianBase          string
	DebianIptables      string
}

// GetAllImages will try to get all images of kube build need.
func GetAllImages() KubeBuildImages {
	k := KubeBuildImages{}
	k.getKubeCross()
	k.getKubePause()
	// DebiamHyperKubeBase not needed
	//	k.getDebianHyperKubeBase()
	k.getDebianBase()
	k.getDebianIptables()
	return k
}

func (k *KubeBuildImages) getKubeCross() {
	cmd := exec.Command("cat", "build/build-image/cross/VERSION")
	out, _ := cmd.CombinedOutput()
	version := strings.TrimSpace(string(out))
	log.Printf("KubeCross version %s", version)
	kubeCross := "us.gcr.io/k8s-artifacts-prod/build-image/kube-cross:" + version
	k.KubeCross = kubeCross
}

func (k *KubeBuildImages) getKubePause() {
	cmd := exec.Command("grep", "-E", "TAG=|TAG =", "build/pause/Makefile")
	out, _ := cmd.CombinedOutput()
	stringOut := string(out)
	outs := strings.Split(stringOut, "=")
	if len(outs) > 1 {
		version := outs[1]
		version = strings.TrimSpace(version)
		log.Printf("KubePause version %s", version)
		kubePause := "k8s.gcr.io/pause-amd64:" + version
		k.KubePause = kubePause
	}
}

func (k *KubeBuildImages) getDebianHyperKubeBase() {
	cmd := exec.Command("grep", "-E", "TAG=|TAG =", "build/debian-hyperkube-base/Makefile")
	out, _ := cmd.CombinedOutput()
	stringOut := string(out)
	outs := strings.Split(stringOut, "=")
	if len(outs) > 1 {
		version := outs[1]
		version = strings.TrimSpace(version)
		log.Printf("DebianHyperKubeBase version %s", version)
		debianHyperKubeBase := "us.gcr.io/k8s-artifacts-prod/build-image/debian-hyperkube-base-amd64:" + version
		k.DebianIptables = debianHyperKubeBase
	}
}

func (k *KubeBuildImages) getDebianBase() {
	cmd := exec.Command("grep", "debian_base_version=", "build/common.sh")
	out, _ := cmd.CombinedOutput()
	stringOut := string(out)
	outs := strings.Split(stringOut, "=")
	if len(outs) > 1 {
		version := outs[1]
		version = strings.TrimSpace(version)
		log.Printf("DebianBase version %s", version)
		debianBase := "us.gcr.io/k8s-artifacts-prod/build-image/debian-base-amd64:" + version
		k.DebianBase = debianBase

	}
}

func (k *KubeBuildImages) getDebianIptables() {
	cmd := exec.Command("grep", "debian_iptables_version=", "build/common.sh")
	out, _ := cmd.CombinedOutput()
	stringOut := string(out)
	outs := strings.Split(stringOut, "=")
	if len(outs) > 1 {
		version := outs[1]
		version = strings.TrimSpace(version)
		log.Printf("DebianIptables version %s", version)
		debianIptables := "us.gcr.io/k8s-artifacts-prod/build-image/debian-iptables-amd64:" + version
		k.DebianIptables = debianIptables
	}
}
