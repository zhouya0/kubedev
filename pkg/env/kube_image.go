package env

import (
	"log"
	"os/exec"
	"strings"
)

// RE_PULL_REMOTE_IMG=false
// if [ "$RE_PULL_REMOTE_IMG" == "true" ]; then
// echo "================================================"
// echo "get remote images versions"
// VERSION_CROSS=$(cat build/build-image/cross/VERSION)
// echo 'kube-cross version:' ${VERSION_CROSS}
// VERSION_PAUSE=$(grep -E "TAG=|TAG =" build/pause/Makefile | cut -d "=" -f 2 | awk '$1=$1')
// echo 'kube-pause version:' ${VERSION_PAUSE}
// VERSION_DEBIAN_HYPERKUBE_BASE=$(grep -E "TAG=|TAG =" build/debian-hyperkube-base/Makefile | cut -d "=" -f 2 | awk '$1=$1')
// echo 'debian-hyperkube version:' ${VERSION_DEBIAN_HYPERKUBE_BASE}
// #VERSION_DEBIAN_BASE=$( grep "debian_base_version=" build/common.sh | cut -d "=" -f 2 | awk '$1=$1')
// VERSION_DEBIAN_BASE=$(grep "TAG ?=" build/debian-base/Makefile | cut -d "=" -f 2 | awk '$1=$1')
// echo 'debian-base version:' ${VERSION_DEBIAN_BASE}
// VERSION_DEBIAN_IPTABLES=$(grep "debian_iptables_version=" build/common.sh | cut -d "=" -f 2 |awk '$1=$1')
// echo 'debian-iptables version:' ${VERSION_DEBIAN_IPTABLES}
// docker_wrapper pull k8s.gcr.io/kube-cross:${VERSION_CROSS}
// docker_wrapper pull k8s.gcr.io/pause-amd64:${VERSION_PAUSE}
// docker_wrapper pull k8s.gcr.io/debian-base-amd64:${VERSION_DEBIAN_BASE}
// docker_wrapper pull k8s.gcr.io/debian-iptables-amd64:${VERSION_DEBIAN_IPTABLES}
// docker_wrapper pull k8s.gcr.io/debian-hyperkube-base-amd64:${VERSION_DEBIAN_HYPERKUBE_BASE}
// fi

type KubeBuildImages struct {
	KubeCross           string
	KubePause           string
	DebianHyperKubeBase string
	DebianBase          string
	DebianIptables      string
}

func (k *KubeBuildImages) getKubeCross() {
	cmd := exec.Command("cat", "build/build-image/cross/VERSION")
	out, _ := cmd.CombinedOutput()
	log.Printf("KubeCross version %s", string(out))
	k.KubeCross = string(out)
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
		k.KubePause = version
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
		k.DebianIptables = version
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
		k.DebianBase = version
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
		k.DebianIptables = version
	}
}
