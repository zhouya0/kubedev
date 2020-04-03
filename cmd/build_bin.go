package cmd

import (
	"kubedev/pkg/build/bin"

	"github.com/spf13/cobra"
)

var supportBinaryName []string = []string{"kubelet", "kubectl"}

var binCmd = &cobra.Command{
	Use:   "bin",
	Short: "Build binaries for kubernetes components",
	Long:  "Build binaries for kubernetes components",

	Run: func(cmd *cobra.Command, args []string) {
		arch, _ := cmd.Flags().GetString("arch")
		if arch == "" {
			arch = "linux/amd64"
		}
		err := BuildBinaryComponents(args, arch)
		if err != nil {
			Fatal(err.Error(), DefaultErrorExitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(binCmd)
	binCmd.Flags().StringP("arch", "a", "", "The binary build arch, could be linux/amd64 or linux/arm64")
}

func BuildBinaryComponents(args []string, arch string) error {
	return bin.BuildBinary(args, arch)
}
