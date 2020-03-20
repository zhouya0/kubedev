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

type KubeVersion struct {
	KubeGitCommit    string
	KubeGitTreeState string
	KubeGitVersion   string
	KubeGitMajor     string
	KubeGitMinor     string
}

func (k *KubeVersion) getKubeGitCommit() error {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Git commit is: %v", string(out))
	k.KubeGitCommit = string(out)
	return nil
}

func (k *KubeVersion) getKubeGitVersion() error {
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
	log.Printf("Kube version is: %v", cutVersion)
	k.KubeGitVersion = cutVersion[0] + cutVersion[1]
	return nil
}

func NewKubeVersion() KubeVersion {
	k := KubeVersion{}
	k.getKubeGitCommit()
	k.getKubeGitVersion()
	return k
}
