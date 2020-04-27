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
		err := BuildBinaryComponents(args, arch)
		if err != nil {
			Fatal(err.Error(), DefaultErrorExitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(binCmd)
	AddArchFlag(binCmd, "linux/amd64")
}

func BuildBinaryComponents(args []string, arch string) error {
	return bin.BuildBinary(args, arch)
}
