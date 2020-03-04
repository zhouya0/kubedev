package main

import (
	"fmt"
	"kubedev/pkg/image"
)

func main() {
	fmt.Println("test")
	err := image.PullImage("k8s.gcr.io/pause-amd64 ")
	if err != nil {
		fmt.Println(err.Error())
	}
}
