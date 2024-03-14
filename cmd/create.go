package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aborroy/alf-k8s/pkg"
	"github.com/spf13/cobra"

	"eagain.net/go/ntlmv2hash"
)

// Default values
const OutputRootPath string = "output"
const TemplateRootPath string = "templates"
const KubernetesEngine string = "docker-desktop"
const DefaultAdminPassword string = "209c6174da490caeb422f3fa5a7ae634" // admin

// Parameters mapping
var interactive bool
var version string
var outputDirectory string
var kubernetes string
var tls bool
var adminPass string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create assets to deploy Alfresco in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {

		var outputRoot = OutputRootPath
		if outputDirectory != "" {
			outputRoot = outputDirectory
		}

		values := pkg.Parameters{}
		if interactive {
			values = pkg.GetPromptValues()
			values.Secret = pkg.GenerateRandomString(24)
			values.AdminPassword = strings.ToLower(ntlmv2hash.NTPasswordHash(values.AdminPassword))
		} else {
			var kubernetesEngine = KubernetesEngine
			if kubernetes != "" {
				kubernetesEngine = kubernetes
			}
			var tlsEnabled = false
			if tls {
				tlsEnabled = true
			}
			var adminPassword = DefaultAdminPassword
			if adminPass != "" {
				adminPassword = strings.ToLower(ntlmv2hash.NTPasswordHash(adminPass))
			}
			values = pkg.Parameters{
				Version:       version,
				Secret:        pkg.GenerateRandomString(24),
				Kubernetes:    kubernetesEngine,
				TLS:           tlsEnabled,
				AdminPassword: adminPassword,
			}
		}

		templateList, err := pkg.EmbedWalk("templates")
		if err != nil {
			panic(err)
		}

		os.MkdirAll(outputRoot, os.ModePerm)

		for _, t := range templateList {
			outputFile := outputRoot + "/" + t
			position := strings.Index(outputFile, "/templates/")
			outputFile = outputFile[position+len("/templates/"):]
			os.MkdirAll(filepath.Dir(outputRoot+"/"+outputFile), os.ModePerm)
			f, _ := os.Create(outputRoot + "/" + outputFile)
			name := filepath.Base(t)
			tpl, _ := template.New(name).ParseFS(pkg.TemplateFs, t)
			err = tpl.ExecuteTemplate(f, name, values)
			if err != nil {
				panic(err)
			}
			if path.Ext(f.Name()) == ".sh" {
				f.Chmod(0755)
			}
			err = pkg.VerifyOutputFile(f.Name())
			if err != nil {
				panic(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Input values replying to command line prompts instead of using command line parameters")
	createCmd.Flags().StringVarP(&version, "version", "v", "", "Version of ACS to be deployed (23.1 or 23.2)")
	createCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "Local Directory to write produced assets, 'output' by default")
	createCmd.Flags().StringVarP(&kubernetes, "kubernetes", "k", "", "Kubernetes cluster: docker-desktop (default) or kind")
	createCmd.Flags().BoolVarP(&tls, "tls", "t", false, "Enable TLS protocol for ingress")
	createCmd.Flags().StringVarP(&adminPass, "password", "p", "", "Password for admin user")
	createCmd.MarkFlagsMutuallyExclusive("interactive", "version")
	createCmd.MarkFlagsMutuallyExclusive("interactive", "kubernetes")
	createCmd.MarkFlagsMutuallyExclusive("interactive", "tls")
	createCmd.MarkFlagsMutuallyExclusive("interactive", "password")
}
