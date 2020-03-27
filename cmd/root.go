package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"kubedev/pkg/env"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// cfgFile is the path of config file you can give.
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kubedev",
	Short: "A development tool for kubernetes",
	Long:  `A development tool for kubernetes`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubedev.yaml)")

	// Get this clear
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		// can't have postprefix here
		viper.SetConfigName(".kubedev")
	}

	// read in environment variables that match
	viper.AutomaticEnv()
	fmt.Println("I'm excuted")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		viper.Unmarshal(&env.Config)
	} else if err != nil {
		fmt.Println(err.Error())
	}
}
