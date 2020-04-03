package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

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

func CheckExist(path string) bool {
	//	str, _ := os.Getwd()
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CopyFile(scrFile, destFile string) error {
	file, err := os.Open(scrFile)
	if err != nil {
		return err
	}
	defer file.Close()

	destdir := path.Dir(destFile)
	if !CheckExist(destdir) {
		err := os.Mkdir(destdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	io.Copy(dest, file)
	return nil
}

func WriteFile(name string, filedata string) error {
	data := []byte(filedata)
	err := ioutil.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
