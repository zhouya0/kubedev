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
		err := BuildBinaryComponents(args)
		if err != nil {
			Fatal(err.Error(), DefaultErrorExitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(binCmd)
}

func BuildBinaryComponents(args []string) error {
	return bin.BuildBinary(args)
}
