package env

// KubeDevConfig is the only config file for kubedev
type KubeDevConfig struct {
	DockerRegistry string
	DockerTag      string
	KubeVersion    string
	BuildPlatform  string
}

var Config KubeDevConfig

var (
	BuildIcon   string = "🔨"
	ImageIcon   string = "💽"
	WriteIcon   string = "📝"
	PackageIcon string = "📦"
)
