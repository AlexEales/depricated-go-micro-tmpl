package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/AlexEales/go-micro-tmpl/tools/heph/helm"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls the projects microservices",
	Long:  `Uninstalls the projects microservices from the k8s cluster currently active in the k8s config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := uninstall(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func uninstall() error {
	helmClient, err := helm.NewClient()
	if err != nil {
		return err
	}

	charts, err := helmClient.ListCharts()
	if err != nil {
		return err
	}

	log.Infof("uninstalling helm charts: {%s}", strings.Join(charts, ", "))
	for _, chart := range charts {
		if err := helmClient.UninstallChart(chart); err != nil {
			return err
		}
		log.Infof("uninstalled chart %s", chart)
	}

	log.Infoln("uninstalling all k8s resources")
	if err := exec.Command("kubectl", "delete", "all", "--all").Run(); err != nil {
		return err
	}

	log.Infoln("cleared local cluster")
	return nil
}
