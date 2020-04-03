package util

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

func GetHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return home
}

func CheckfileExist(path string) bool {
	str, _ := os.Getwd()
	_, err := os.Stat(str + "/" + path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
