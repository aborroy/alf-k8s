package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
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
var version string
var outputDirectory string
var kubernetes string
var tls string
var adminPass string

// Template Variables
type Values struct {
	Version       string // ACS Version (23.2, 23.1...)
	Secret        string // Shared secret string for Repo and Solr communication
	Kubernetes    string // Kubernetes cluster (Docker Desktop, KinD)
	TLS           bool   // Enable TLS in ingress controller for https access
	AdminPassword string // Password for admin user ('admin' is default password)
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
		var adminPassword = DefaultAdminPassword
		if adminPass != "" {
			// Alfresco accepts only lower case NTLM passwords
			adminPassword = strings.ToLower(ntlmv2hash.NTPasswordHash(adminPass))
		}

		values := Values{
			version,
			pkg.GenerateRandomString(24),
			kubernetesEngine,
			tlsEnabled,
			adminPassword}

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
	createCmd.Flags().StringVarP(&adminPass, "password", "p", "", "Password for admin user")
	createCmd.MarkFlagRequired("version")
}
