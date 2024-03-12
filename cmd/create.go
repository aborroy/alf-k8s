package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aborroy/alf-k8s/pkg"
	"github.com/spf13/cobra"
)

// Default values
const OutputRootPath string = "output"
const TemplateRootPath string = "templates"

// Parameters mapping
var version string
var outputDirectory string

// Template Variables
type Values struct {
	Version string
	Secret  string
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create assets to deploy Alfresco in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {

		var outputRoot = OutputRootPath
		if outputDirectory != "" {
			outputRoot = outputDirectory
		}

		values := Values{version, pkg.GenerateRandomString(24)}

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
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&version, "version", "v", "", "Version of ACS to be deployed (23.1 or 23.2)")
	createCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "Local Directory to write produced assets")
	createCmd.MarkFlagRequired("version")
}
