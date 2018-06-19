package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skrt",
	Short: "skrt (Simple Kube Reliability Toolkit) is a simple tool providing simple and opinionated SRE infrastructure based on Kubernetes.",
	Long: `skrt (Simple Kube Reliability Toolkit) is a simple tool providing simple and opinionated SRE infrastructure based on Kubernetes.
                
See https://github.com/hekonsek/skrt`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}