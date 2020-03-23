package env

// echo "================================================"
// echo "get kubernetes version"
// KUBE_GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null)
// KUBE_GIT_VERSION=$(git describe --tags --abbrev=14 "${KUBE_GIT_COMMIT}" | cut -d "-" -f 1,2)
// #KUBE_VERSION=${KUBE_GIT_VERSION##v} # remove the leading v, 将v1.15.3改为1.15.3
// KUBE_VERSION=${KUBE_GIT_VERSION}
// echo "================================================"
// echo "create kubernetes version file"
// cat <<EOF >"dce_version"
// KUBE_GIT_COMMIT=${KUBE_GIT_COMMIT-}
// KUBE_GIT_TREE_STATE='clean'
// KUBE_GIT_VERSION=${KUBE_GIT_VERSION-}
// KUBE_GIT_MAJOR='${KUBE_VERSION:0:1}'
// KUBE_GIT_MINOR='${KUBE_VERSION:2:2}'
// EOF

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// cmd := exec.Command("cat", "build/build-image/cross/VERSION")
// out, _ := cmd.CombinedOutput()

var cleanTreeState string = "clean"

type KubeVersion struct {
	KubeGitCommit    string
	KubeGitTreeState string
	KubeGitVersion   string
	KubeGitMajor     string
	KubeGitMinor     string
}

func (k *KubeVersion) setKubeGitCommit(s string) error {
	if s != "" {
		k.KubeGitCommit = s
		return nil
	}
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	k.KubeGitCommit = string(out)
	log.Printf("Git commit is: %v \n", k.KubeGitCommit)
	return nil
}

func (k *KubeVersion) setKubeGitVersion(s string) error {
	if s != "" {
		k.KubeGitVersion = s
		return nil
	}
	if k.KubeGitCommit == "" {
		return fmt.Errorf("Can't get kube git version!")
	}
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=14", strings.TrimSpace(k.KubeGitCommit))
	fullVersion, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	cutVersion := strings.Split(string(fullVersion), "-")
	k.KubeGitVersion = cutVersion[0] + cutVersion[1]
	log.Printf("Kube version is: %v \n", k.KubeGitVersion)
	return nil
}

func (k *KubeVersion) setKubeTreeState(s string) error {
	if s != "" {
		k.KubeGitTreeState = s
		return nil
	}
	k.KubeGitTreeState = cleanTreeState
	log.Printf("Git tree state is: %v \n", k.KubeGitTreeState)
	return nil
}

func (k *KubeVersion) setKubeGitMajorAndMinor() error {
	if k.KubeGitVersion == "" {
		return fmt.Errorf("Can't get kube major and minor version!")
	}
	versions := strings.Split(k.KubeGitVersion, ".")

	k.KubeGitMajor = strings.Split(versions[0], "v")[1]
	k.KubeGitMinor = versions[1]
	log.Printf("Kube git major is: %v \n", k.KubeGitMajor)
	log.Printf("Kube git minor is: %v \n", k.KubeGitMinor)
	return nil
}

func NewKubeVersion() KubeVersion {
	k := KubeVersion{}
	k.setKubeGitCommit("")
	k.setKubeGitVersion("")
	k.setKubeTreeState("")
	k.setKubeGitMajorAndMinor()
	return k
}
