package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "heph",
	Short: "Heph is a basic set of utilities for automating the deployment of microservices",
}

func init() {
	rootCmd.AddCommand(deployCmd, installCmd, uninstallCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
