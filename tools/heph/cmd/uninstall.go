package cmd

import "github.com/spf13/cobra"

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the projects microservices",
	Long:  `Uninstalls the projects microservices from the k8s cluster currently active in the k8s config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Uninstalling project")
	},
}
