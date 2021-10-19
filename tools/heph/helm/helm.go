package helm

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// TODO: Would be good to test this but would rely on being able to either:
// 		 - mock the exec.Command calls
// 		 - have a mock helm binary run in place of normal helm

var log = logrus.New()

// Chart represents a helm chart
type Chart struct {
	Name          string
	OverridesFile string
}

// Client defines a interface to interact with helm repositories and charts
type Client interface {
	AddRepository(string, string) error
	AddRepositories(map[string]string) error
	InstallChart(string, *Chart) error
	InstallCharts(map[string]*Chart) error
	ListCharts() ([]string, error)
	ListRepositories() ([]string, error)
	RemoveRepository(string) error
	UninstallChart(string) error
}

// NewClient returns a new helm client, erroring if helm is not installed
func NewClient() (Client, error) {
	_, err := exec.LookPath("helm")
	if err != nil {
		err = fmt.Errorf("error creating helm client, helm does not exist in path")
		log.Errorln(err)
		return nil, err
	}
	return &helmClient{}, nil
}

type helmClient struct{}

// AddRepository adds a repository at the provided url and assigns it the given name
func (c *helmClient) AddRepository(name string, url string) error {
	out, err := exec.Command("helm", "repo", "add", name, url).CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"name": name,
			"url":  url,
		}).Errorf("error occurred while adding helm repository")
		log.Println(string(out))
		return err
	}

	log.Infof("added helm repository %s <%s>", name, url)
	return nil
}

// AddRepositories adds all the specified repositories in the map (name: url)
func (c *helmClient) AddRepositories(repos map[string]string) error {
	for name, url := range repos {
		if err := c.AddRepository(name, url); err != nil {
			return err
		}
	}

	return nil
}

// InstallChart installs a specifed chart with the given name on the currently selected
// kubernetes context.
func (c *helmClient) InstallChart(name string, chart *Chart) error {
	if chart.OverridesFile == "" {
		return installChart(name, chart)
	} else {
		return installChartWithOverrides(name, chart)
	}
}

// InstallCharts installs all the specified charts in the map (name: chart)
func (c *helmClient) InstallCharts(charts map[string]*Chart) error {
	for name, chart := range charts {
		if err := c.InstallChart(name, chart); err != nil {
			return err
		}
	}

	return nil
}

func installChart(name string, chart *Chart) error {
	out, err := exec.Command("helm", "install", name, chart.Name).CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"name": chart.Name,
		}).Errorf("error occurred while installing helm chart")
		log.Println(string(out))
		return err
	}

	log.Infof("installed helm repository %s", chart.Name)
	return nil
}

func installChartWithOverrides(name string, chart *Chart) error {
	out, err := exec.Command("helm", "install", name, chart.Name, "--values", chart.OverridesFile).CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"name":           chart.Name,
			"overrides_file": chart.OverridesFile,
		}).Errorf("error occurred while installing helm chart")
		log.Println(string(out))
		return err
	}

	log.Infof("installed helm repository %s with override vaules %s", chart.Name, chart.OverridesFile)
	return nil
}

// ListCharts lists all the installed chart names
func (c *helmClient) ListCharts() ([]string, error) {
	out, err := exec.Command("helm", "list", "-q").CombinedOutput()
	if err != nil {
		log.WithError(err).
			Error("error occured while listing installed charts")
		return nil, err
	}

	formattedOutput := strings.TrimSpace(string(out))
	return strings.Split(formattedOutput, "\n"), nil
}

// ListCharts lists all the installed repository names
func (c *helmClient) ListRepositories() ([]string, error) {
	out, err := exec.Command("helm", "repo", "list", "-o", "json").CombinedOutput()
	if err != nil {
		log.WithError(err).
			Error("error occured while listing installed repositories")
		return nil, err
	}

	repositories := []*struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{}
	if err := json.Unmarshal(out, &repositories); err != nil {
		log.WithError(err).Error("error parsing list repositories output")
		return nil, err
	}

	repositoryNames := make([]string, len(repositories))
	for i, repo := range repositories {
		repositoryNames[i] = repo.Name
	}

	return repositoryNames, nil
}

// RemoveRepository removes the repository with the provided name
func (c *helmClient) RemoveRepository(name string) error {
	out, err := exec.Command("helm", "repo", "remove", name).CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"name": name,
		}).Errorf("error occurred while removing helm repository")
		log.Println(string(out))
		return err
	}

	log.Infof("removed helm repository %s", name)
	return nil
}

// UninstallChart uninstalls the resources deployed by the specified chart
// from the currently selected kubernetes context
func (c *helmClient) UninstallChart(name string) error {
	out, err := exec.Command("helm", "uninstall", name).CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"name": name,
		}).Errorf("error occurred while uninstalling helm chart")
		log.Println(string(out))
		return err
	}

	log.Infof("uninstalled helm chart %s", name)
	return nil
}
