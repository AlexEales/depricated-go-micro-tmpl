package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path to service]",
	Short: "Deploys the specified service",
	Long:  `Deploys the specified service to the k8s cluster currently active in the kube config.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Deploying service %s", args[0])

		if err := validateServiceDir(args[0]); err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if err := deploy(args[0]); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	},
}

func deploy(dir string) error {
	inputPath := fmt.Sprintf("%s/main.go", dir)
	outputPath := fmt.Sprintf("%s/docker/%s.bin", dir, path.Base(dir))
	if err := exec.Command("go", "build", "-o", outputPath, inputPath).Run(); err != nil {
		return err
	}
	log.Infof("built service binary %s", outputPath)

	sourcePath := fmt.Sprintf("./%s/docker", dir)
	// TODO: Add this as some kind of global variable so we can avoid find-replace
	tagName := fmt.Sprintf("go-micro-tmpl/%s", path.Base(dir))
	if err := exec.Command("docker", "build", "-t", tagName, sourcePath).Run(); err != nil {
		return err
	}
	log.Infof("built docker image %s", tagName)

	resourcesPath := fmt.Sprintf("%s/k8s", dir)
	if err := exec.Command("kubectl", "apply", "-f", resourcesPath).Run(); err != nil {
		return err
	}
	log.Infof("deployed k8s resources from %s", resourcesPath)

	path := fmt.Sprintf("%s/docker/%s.bin", dir, path.Base(dir))
	if err := os.Remove(path); err != nil {
		return err
	}
	log.Infof("service %s deployed", dir)

	return nil
}

func validateServiceDir(dir string) error {
	if _, err := os.Stat(dir + "/docker/Dockerfile"); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile missing")
	}

	if _, err := os.Stat(dir + "/k8s"); os.IsNotExist(err) {
		return fmt.Errorf("k8s dir missing")
	}

	if _, err := os.Stat(dir + "/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("no main.go file present")
	}

	return nil
}
