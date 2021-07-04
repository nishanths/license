package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

<<<<<<< HEAD
<<<<<<< HEAD
var licenses = map[string]struct {
=======
var knownLicenses = map[string]struct {
>>>>>>> Duplicate licensesList from list.go
=======
var licenses = map[string]struct {
>>>>>>> Rename knownLicenses to simply licenses
	longName string
	template string
}{
	"agpl-3.0":     {"GNU Affero General Public License v3.0", Agpl30Template},
	"apache-2.0":   {"Apache License 2.0", Apache20Template},
	"bsd-2-clause": {"BSD 2-Clause \"Simplified\" License", Bsd2ClauseTemplate},
	"bsd-3-clause": {"BSD 3-Clause \"New\" or \"Revised\" License", Bsd3ClauseTemplate},
	"cc0-1.0":      {"Creative Commons Zero v1.0 Universal", Cc010Template},
	"epl-2.0":      {"Eclipse Public License 2.0", Epl20Template},
	"free-art-1.3": {"Free Art License 1.3", FreeArt13Template},
	"gpl-2.0":      {"GNU General Public License v2.0", Gpl20Template},
	"gpl-3.0":      {"GNU General Public License v3.0", Gpl30Template},
	"lgpl-2.1":     {"GNU Lesser General Public License v2.1", Lgpl21Template},
	"lgpl-3.0":     {"GNU Lesser General Public License v3.0", Lgpl30Template},
	"mit":          {"MIT License", MitTemplate},
	"mpl-2.0":      {"Mozilla Public License 2.0", Mpl20Template},
	"unlicense":    {"The Unlicense", UnlicenseTemplate},
	"wtfpl":        {"Do What The Fuck You Want To Public License", WtfplTemplate},
}

func printList() {
<<<<<<< HEAD
<<<<<<< HEAD
	keys := make([]string, 0, len(licenses))

	for key := range licenses {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		stdout.Printf("%-14s(%s)", key, licenses[key].longName)
=======
	for key, license := range knownLicenses {
=======
	for key, license := range licenses {
>>>>>>> Rename knownLicenses to simply licenses
		stdout.Printf("%-14s(%s)", key, license.longName)
>>>>>>> Duplicate licensesList from list.go
	}
}

func printLicense(license, output, name, year string) {
<<<<<<< HEAD
<<<<<<< HEAD
	file, ok := licenses[license]
=======
	licenseData, ok := knownLicenses[license]
>>>>>>> Duplicate licensesList from list.go
=======
	licenseData, ok := licenses[license]
>>>>>>> Rename knownLicenses to simply licenses
	if !ok {
		stderr.Printf("unknown license %q\nrun \"license -list\" for list of available licenses", license)
		os.Exit(2)
	}

	t, err := template.New("license").Parse(file.template)
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
