package main

import (
	"embed"

	"github.com/aborroy/alf-k8s/cmd"
	"github.com/aborroy/alf-k8s/pkg"
)

//go:embed all:templates
var templateFs embed.FS

func main() {
	pkg.TemplateFs = templateFs
	cmd.TemplateFs = templateFs
	cmd.Execute()
}
