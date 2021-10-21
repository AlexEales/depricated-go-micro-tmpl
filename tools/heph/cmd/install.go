package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the projects microservices and sets up the project",
	Long:  `Installs the projects microservices and sets up the project on the k8s cluster currently active in the k8s config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Installing project")
	},
}
