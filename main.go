package main

import (
	"fmt"
	"kubedev/pkg/env"
)

func main() {
	err := env.WriteVersionFile("dce_version")
	if err != nil {
		fmt.Println(err.Error())
	}
}
