package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/AlexEales/go-micro-tmpl/tools/heph/helm"
	"github.com/AlexEales/go-micro-tmpl/tools/heph/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var log = logrus.New()

func main() {
	helmClient, err := helm.NewClient()
	if err != nil {
		log.Panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		log.Panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
	k8sClient := k8s.NewClient(clientset)

	if err := helmClient.AddRepositories(map[string]string{
		"bitnami":   "https://charts.bitnami.com/bitnami",
		"hashicorp": "https://helm.releases.hashicorp.com",
	}); err != nil {
		log.Panic(err)
	}

	if err := helmClient.InstallCharts(map[string]*helm.Chart{
		"postgres": {
			Name:          "bitnami/postgresql",
			OverridesFile: "infra/postgres/override-values.yaml",
		},
	}); err != nil {
		log.Panic(err)
	}

	if err := k8sClient.WaitForPodsToBeReady(context.TODO(), "default", []string{"postgres-primary-0", "postgres-read-0"}, time.Minute); err != nil {
		log.Panic(err)
	}

	log.Infoln("Installed!")
}
