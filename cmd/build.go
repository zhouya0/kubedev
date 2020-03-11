package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the kubernetes components",
	Long:  "Build the kubernetes components",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build called")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
