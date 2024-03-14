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
}

var qs = []*survey.Question{
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
			Message: "Choose the password for your admin user",
			Default: "admin",
		},
	},
}

func GetPromptValues() Parameters {
	answers := Parameters{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		panic(err)
	}
	return answers
}
