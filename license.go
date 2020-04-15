package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

//go:generate go-bindata -o templates.go .templates/

func printLicense(license, output, name, year string) {
	data, err := Asset(fmt.Sprintf(".templates/%s.tmpl", license))
	if err != nil {
		stderr.Printf("unknown license %q\nrun \"license -list\" for list of available licenses", license)
		os.Exit(2)
	}

	t, err := template.New("license").Parse(string(data))
	if err != nil {
		stderr.Printf("internal: failed to parse license template for %s", license)
		os.Exit(1)
	}

	var outFile io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(filepath.Clean(output))
		if err != nil {
			stderr.Printf("failed to create file %s: %s", output, err)
			os.Exit(1)
		}
		outFile = f
	}

	if err := t.Execute(outFile, struct {
		Name string
		Year string
	}{name, year}); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}