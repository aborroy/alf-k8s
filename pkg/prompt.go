package pkg

import (
	"github.com/AlecAivazis/survey/v2"
)

type Parameters struct {
	Version       string // ACS Version (23.2, 23.1...)
	Secret        string // Shared secret string for Repo and Solr communication
	Kubernetes    string // Kubernetes cluster (Docker Desktop, KinD)
	TLS           bool   // Enable TLS in ingress controller for https access
	AdminPassword string // Password for admin user ('admin' is default password)
	DockerUser    string // Username for Docker Hub
	DockerPass    string // Password of the username for Docker Hub
	DockerAuth    string // Base64 encode for Docker Hub credentials
}

var basicQuestions = []*survey.Question{
	{
		Name: "version",
		Prompt: &survey.Select{
			Message: "Which ACS version do you want to use?",
			Options: []string{"23.2", "23.1"},
			Default: "23.2",
		},
	},
	{
		Name: "kubernetes",
		Prompt: &survey.Select{
			Message: "What Kubernetes cluster do you want to use?",
			Options: []string{"docker-desktop", "kind"},
			Default: "docker-desktop",
		},
	},
	{
		Name: "tls",
		Prompt: &survey.Confirm{
			Message: "Do you want to use HTTPs for Ingress?",
		},
	},
	{
		Name: "adminPassword",
		Prompt: &survey.Input{
			Message: "Choose the password for your ACS admin user",
			Default: "admin",
		},
	},
}

var kindQuestions = []*survey.Question{
	{
		Name: "dockerUser",
		Prompt: &survey.Input{
			Message: "Provide an existing username in Docker Hub",
		},
	},
	{
		Name: "dockerPass",
		Prompt: &survey.Input{
			Message: "Provide a password for the username in Docker Hub",
		},
	},
}

func GetPromptValues() Parameters {
	answers := Parameters{}
	err := survey.Ask(basicQuestions, &answers)
	if err != nil {
		panic(err)
	}
	if answers.Kubernetes == "kind" {
		err := survey.Ask(kindQuestions, &answers)
		if err != nil {
			panic(err)
		}
	}
	return answers
}
