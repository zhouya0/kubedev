package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the kubernetes components",
	Long:  "Build the kubernetes components",

	Run: DefaultSubCommandRun(os.Stderr),
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(imageCmd)
}
