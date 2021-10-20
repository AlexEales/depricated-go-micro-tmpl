package main

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/AlexEales/go-micro-tmpl/tools/heph/helm"
	"github.com/AlexEales/go-micro-tmpl/tools/heph/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Command func() error

var (
	commandMap = map[string]Command{
		"install":   install,
		"uninstall": uninstall,
	}
	// TODO: not a fan of the logrus formatting should make a common lib with a common format
	log = logrus.New()
)

func install() error {
	helmClient, err := helm.NewClient()
	if err != nil {
		return err
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k8sClient := k8s.NewClient(clientset)

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

	if err := k8sClient.WaitForPodsToBeReady(context.TODO(), "default", []string{"postgres-primary-0", "postgres-read-0"}, time.Minute); err != nil {
		return err
	}

	log.Infoln("Installed!")
	return nil
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

	for _, chart := range charts {
		if err := helmClient.UninstallChart(chart); err != nil {
			return err
		}
		log.Infof("uninstalled chart %s", chart)
	}

	return nil
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Panic("no arguments provided")
	}

	cmd, ok := commandMap[args[0]]
	if !ok {
		log.Panicf("command <%s> not known", args[0])
	}

	if err := cmd(); err != nil {
		log.Panic(err)
	}
}
