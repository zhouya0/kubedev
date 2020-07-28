package env

// KubeDevConfig is the only config file for kubedev
type KubeDevConfig struct {
	DockerRegistry      string
	DockerTag           string
	OverrideKubeVersion string
	BuildPlatform       string
	FastBuild string
}

var Config KubeDevConfig

var (
	BuildIcon   string = "🔨"
	ImageIcon   string = "💽"
	WriteIcon   string = "📝"
	PackageIcon string = "📦"
)
