package env

// KubeDevConfig is the only config file for kubedev
type KubeDevConfig struct {
	DockerRegistry string
	DockerTag      string
}


var Config KubeDevConfig