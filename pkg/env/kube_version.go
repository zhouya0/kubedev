package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

const (
	KubeBinPath string = "_output/dockerized/bin"
)

var KubeVersionFile string = "kube_version"

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
	if len(cutVersion) > 1 {
		k.KubeGitVersion = strings.Join([]string{cutVersion[0], cutVersion[1]}, "-")
	} else {
		k.KubeGitVersion = cutVersion[0]
	}
	return nil
}

func (k *KubeVersion) setKubeTreeState(s string) error {
	if s != "" {
		k.KubeGitTreeState = s
		return nil
	}
	k.KubeGitTreeState = cleanTreeState
	return nil
}

func (k *KubeVersion) setKubeGitMajorAndMinor() error {
	if k.KubeGitVersion == "" {
		return fmt.Errorf("Can't get kube major and minor version!")
	}
	versions := strings.Split(k.KubeGitVersion, ".")

	k.KubeGitMajor = strings.Split(versions[0], "v")[1]
	k.KubeGitMinor = versions[1]
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

func generateContent() string {
	var s string
	k := NewKubeVersion()
	s = s + fmt.Sprintf("KUBE_GIT_COMMIT=%s\n", k.KubeGitCommit)
	s = s + fmt.Sprintf("KUBE_GIT_TREE_STATE=%s\n", k.KubeGitTreeState)
	s = s + fmt.Sprintf("KUBE_GIT_VERSION=%s\n", k.KubeGitVersion)
	s = s + fmt.Sprintf("KUBE_GIT_MAJOR=%s\n", k.KubeGitMajor)
	s = s + fmt.Sprintf("KUBE_GIT_MINOR=%s\n", k.KubeGitMinor)
	return s
}

func WriteVersionFile(name string) error {
	content := generateContent()
	data := []byte(content)
	err := ioutil.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetKubeVersion(k KubeVersion) string {
	return k.KubeGitVersion
}

func GetKubeVersionNoV(k KubeVersion) string {
	version := strings.Split(k.KubeGitVersion, "-")[0]
	trimedVersion := strings.Trim(version, "v")
	return trimedVersion
}

func GetComponentDirName(component string, k KubeVersion) string {
	version := strings.Split(k.KubeGitVersion, "-")[0]
	trimedVersion := strings.Trim(version, "v")
	return strings.Join([]string{component, trimedVersion}, "-")
}
