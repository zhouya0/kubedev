package cmd

import (
	"fmt"
	"kubedev/pkg/build/image"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the kubernetes components",
	Long:  "Build the kubernetes components",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		BuildComponents(args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func BuildComponents(args []string) error {
	return image.BuildImages()
}
