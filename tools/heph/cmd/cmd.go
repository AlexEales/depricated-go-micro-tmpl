package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// TODO: not a fan of the logrus formatting should make a common lib with a common format
	log = logrus.New()
)

var rootCmd = &cobra.Command{
	Use:   "heph",
	Short: "Heph is a basic set of utilities for automating the deployment of microservices",
}

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	rootCmd.AddCommand(deployCmd, installCmd, uninstallCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
