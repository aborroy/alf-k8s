package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"strconv"
	"text/template"

	"github.com/aborroy/alf-k8s/pkg"
	"github.com/spf13/cobra"
)

// Default values
const OutputRootPath string = "output"
const TemplateRootPath string = "templates"
const KubernetesEngine string = "docker-desktop"

// Parameters mapping
var version string
var outputDirectory string
var kubernetes string
var tls string

// Template Variables
type Values struct {
	Version string
	Secret  string
	Kubernetes string
	TLS bool
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create assets to deploy Alfresco in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {

		var outputRoot = OutputRootPath
		if outputDirectory != "" {
			outputRoot = outputDirectory
		}

		var kubernetesEngine = KubernetesEngine
		if kubernetes != "" {
			kubernetesEngine = kubernetes
		}
		var tlsEnabled = false
		if tls != "" {
			tlsEnabled, _ = strconv.ParseBool(tls)
		}

		values := Values{
			version, 
			pkg.GenerateRandomString(24), 
			kubernetesEngine, 
			tlsEnabled}

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
	createCmd.Flags().StringVarP(&version, "version", "v", "", "Version of ACS to be deployed (23.1 or 23.2)")
	createCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "Local Directory to write produced assets, 'output' by default")
	createCmd.Flags().StringVarP(&kubernetes, "kubernetes", "k", "", "Kubernetes cluster: docker-desktop (default) or kind")
	createCmd.Flags().StringVarP(&tls, "tls", "t", "", "Enable TLS protocol for ingress")
	createCmd.MarkFlagRequired("version")
}
