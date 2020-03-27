package main

import (
	"fmt"
	"kubedev/cmd"
	"kubedev/pkg/env"
)

func main() {
	cmd.Execute()
	k := env.NewKubeDevConfig()
	fmt.Println(k)
}
