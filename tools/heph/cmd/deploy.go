package cmd

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy [path to service]",
	Short: "Deploys the specified service",
	Long:  `Deploys the specified service to the k8s cluster currently active in the kube config.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Deploying service %s", args[0])
	},
}
