package env

import (
	"fmt"

	"github.com/spf13/viper"
)

type KubeDevConfig struct {
	DockerRegistry string
	DockerTag      string
}

func NewKubeDevConfig() KubeDevConfig {
	var config KubeDevConfig
	viper.Unmarshal(config)
	fmt.Println(config)
	return config
}
