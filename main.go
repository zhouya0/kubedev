package main

import (
	"fmt"
	"kubedev/pkg/env"
)

func main() {
	kubeVersion := env.NewKubeVersion()
	fmt.Printf("final version is %v \n", kubeVersion)
}
