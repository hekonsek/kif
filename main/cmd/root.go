package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kif",
	Short: "kif (Kube Is Fine) is a small tool providing simple and opinionated SRE infrastructure based on Kubernetes.",
	Long: `kif (Kube Is Fine) is a small tool providing simple and opinionated SRE infrastructure based on Kubernetes.
                
See https://github.com/hekonsek/kif`,
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
