package cmd

import (
	"kubedev/pkg/build/rpm"

	"github.com/spf13/cobra"
)

var supportRPM []string = []string{"kubelet"}

var rpmCmd = &cobra.Command{
	Use:   "rpm",
	Short: "Build rpm for kubernetes components",
	Long:  "Build rpm for kubernetes components",
	Run: func(cmd *cobra.Command, args []string) {
		arch, _ := cmd.Flags().GetString("arch")
		if arch == "" {
			arch = "linux/amd64"
		}
		err := BuildRPMComponents(args, arch)
		if err != nil {
			Fatal(err.Error(), DefaultErrorExitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(rpmCmd)
	rpmCmd.Flags().StringP("arch", "a", "", "The RPM build arch, could be linux/amd64 or linux/arm64")
}

func BuildRPMComponents(args []string, arch string) error {
	return rpm.BuildRPM(args, arch)
}
