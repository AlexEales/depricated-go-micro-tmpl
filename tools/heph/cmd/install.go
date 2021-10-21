package cmd

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlexEales/go-micro-tmpl/tools/heph/helm"
	"github.com/AlexEales/go-micro-tmpl/tools/heph/k8s"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the projects microservices and sets up the project",
	Long:  `Installs the projects microservices and sets up the project on the k8s cluster currently active in the k8s config`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := install(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func install() error {
	startTime := time.Now()

	helmClient, err := helm.NewClient()
	if err != nil {
		return err
	}

	k8sClient, err := k8s.NewClient(filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return err
	}

	if err := helmClient.AddRepositories(map[string]string{
		"bitnami":   "https://charts.bitnami.com/bitnami",
		"hashicorp": "https://helm.releases.hashicorp.com",
	}); err != nil {
		return err
	}

	if err := helmClient.InstallCharts(map[string]*helm.Chart{
		"postgres": {
			Name:          "bitnami/postgresql",
			OverridesFile: "infra/postgres/override-values.yaml",
		},
	}); err != nil {
		return err
	}

	pods := []string{
		"postgres-primary-0",
		"postgres-read-0",
	}
	log.Infof("Waiting for pods to be ready: {%s}", strings.Join(pods, ", "))
	if err := k8sClient.WaitForPodsToBeReady(context.TODO(), "default", pods, time.Minute); err != nil {
		return err
	}

	elapsed := time.Since(startTime)
	log.Infof("installed successfully in %.2gm%.2gs", math.Round(elapsed.Minutes()), math.Round(elapsed.Seconds()))
	return nil
}
