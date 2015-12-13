package main

import (
	"encoding/json"
	"flag"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"github.com/nishanths/license/console"
	"github.com/nishanths/license/local"
	"github.com/nishanths/license/remote"
	shutil "github.com/termie/go-shutil"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	tempDirectoryPrefix = "license-"
	applicationVersion  = "0.1.0"
)

func listFailed() {
	console.Error(fmt.Sprintf("license: failed to fetch available licenses: run `$ license bootstrap` before trying again."))
}

func permissionsFailed(p string) {
	console.Error(fmt.Sprintf("Could not access %s. Make sure you have the permissions", p))
}

func createPathFail(p string) {
	console.Error(fmt.Sprintf("Failed to make %s. Make sure you have the permissions", p))
}

func createPathSuccess(p string) {
	console.Info(fmt.Sprintf("Created %s", p))
}

func render(key, fullname string, year int, w io.Writer) bool {
	var c base.Config
	c.Prepare(fullname, "")

	var o base.Option
	o.Name = c.Name
	o.Year = year

	tmpl, err := local.Template(key)
	if err != nil {
		// TODO: error message
		return false
	}

	base.RenderTemplate(tmpl, &o, w)
	return true
}

func generate(args []string) bool {
	if len(args) < 1 {
		// TODO: print error message
		return false
	}

	name := ""
	year := time.Now().Year()
	output := ""
	var f *os.File

	if len(args) > 2 {
		var c base.Config
		c.Prepare("", "")

		generateSet := flag.NewFlagSet("generate", flag.ExitOnError)

		generateSet.StringVar(&name, "n", c.Name, "Specify name to use on license")
		generateSet.StringVar(&name, "name", c.Name, "Specify name to use on license")

		generateSet.IntVar(&year, "y", time.Now().Year(), "Specify year to use of license")
		generateSet.IntVar(&year, "year", time.Now().Year(), "Specify year to use of license")

		generateSet.StringVar(&output, "o", "", "Specify output file name")
		generateSet.StringVar(&output, "output", "", "Specify output file name")

		err := generateSet.Parse(args[1:])
		if err != nil {
			// TODO: print error
			return false
		}
	}

	if output == "" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.Create("./" + output)
		if err != nil {
			// TODO: error message
			return false
		}
	}

	licenseName := strings.ToLower(args[0])
	licenses, err := local.List()

	if err != nil {
		listFailed()
		return false
	}

	for _, l := range licenses {
		if l.Key == licenseName || l.Name == licenseName {
			return render(l.Key, name, year, f)
		}
	}

	return false
}

func update() bool {
	// Bail immediately if we cannot find the user's home directory
	home, err := homedir.Dir()
	if err != nil {
		console.Error("Unable to locate home directory")
		return false
	}

	tempLicensePath, err := ioutil.TempDir(os.TempDir(), tempDirectoryPrefix)
	if err != nil {
		createPathFail("temporary directory")
		return false
	}

	dataPath := path.Join(tempLicensePath, base.DataDirectory)
	rawPath := path.Join(dataPath, base.RawDirectory)
	templatesPath := path.Join(dataPath, base.TemplatesDirectory)
	listFilePath := path.Join(dataPath, base.ListFile)

	defer func() {
		console.Info("Cleaning up...")
		console.Info(fmt.Sprintf("Removing temp directory %s", tempLicensePath))
		os.RemoveAll(tempLicensePath)
		console.Info("Bootstrap complete")
	}()

	// Create data directories
	pathsToMake := []string{rawPath, templatesPath}
	for _, p := range pathsToMake {
		if err := os.MkdirAll(p, 0700); err != nil {
			createPathFail(p)
			return false
		}
	}
	createPathSuccess(fmt.Sprintf("temp directory %s", tempLicensePath))

	// Fetch list from remote
	licenses, err := remote.List()
	if err != nil {
		console.Error("Failed to make licenses list from api.github.com")
		return false
	}
	console.Info("Fetched license list from api.github.com")

	// Serialize the list into JSON
	// Write the serialized JSON to the list file
	serialized, err := json.MarshalIndent(licenses, "", "    ")
	if err != nil {
		console.Error(fmt.Sprintf("Failed to serialize licenses. Please create an issue: %s", base.RepositoryIssuesURL))
		return false
	}

	if err := ioutil.WriteFile(listFilePath, serialized, 0700); err != nil {
		createPathFail(listFilePath)
		return false
	}

	fullLicenses := make([]base.License, 0)

	// Fetch each license's full info
	// - Serialize to JSON and write to disk
	// - Convert to text template and write to disk
	for _, l := range licenses {
		fullLicense, err := remote.Info(&l)
		if err != nil {
			console.Error("Failed to make detailed license info")
			return false
		}

		fullLicenses = append(fullLicenses, fullLicense)

		serialized, err := json.MarshalIndent(fullLicense, "", "    ")
		if err != nil {
			console.Error(fmt.Sprintf("Failed to serialize licenses. Please file an issue: %s", base.RepositoryIssuesURL))
			return false
		}

		rawFilePath := path.Join(rawPath, l.Key+".json")
		if err := ioutil.WriteFile(rawFilePath, serialized, 0700); err != nil {
			createPathFail(rawFilePath)
			return false
		}

		templateData := fullLicense.TextTemplate()
		templateFilePath := path.Join(templatesPath, l.Key+".tmpl")
		if err := ioutil.WriteFile(templateFilePath, []byte(templateData), 0700); err != nil {
			createPathFail(templateFilePath)
			return false
		}
	}

	// Remove exisitng path
	realLicensePath := path.Join(home, base.LicenseDirectory)
	if err := os.RemoveAll(realLicensePath); err != nil && os.IsPermission(err) {
		permissionsFailed(realLicensePath)
		return false
	}

	// Copy temp data to real path
	if err := shutil.CopyTree(tempLicensePath, realLicensePath, nil); err != nil {
		console.Error(fmt.Sprintf("Failed to copy data to %s", realLicensePath))
		return false
	}
	createPathSuccess(fmt.Sprintf("and copied data to %s", realLicensePath))

	return true
}

func version() bool {
	fmt.Println(applicationVersion)
	return true
}

func listHelper(fn func() ([]base.License, error)) bool {
	licenses, err := fn()
	if err != nil {
		listFailed()
		return false
	}

	base.RenderList(&licenses)
	fmt.Println()
	return true
}

func listRemote() bool {
	fmt.Println("Available licenses (remote):\n")
	return listHelper(remote.List)
}

func listLocal() bool {
	fmt.Println("Available licenses (local):\n")
	return listHelper(local.List)
}

type usageLine struct {
	Command     string
	Description string
}

type exampleLine usageLine

func (l *usageLine) String() string {
	return fmt.Sprintf("   %-28s%s", l.Command, l.Description)
}

func (l *exampleLine) String() string {
	return fmt.Sprintf("   %-40s%s", l.Command, l.Description)
}

func help() bool {
	// Heading
	fmt.Println("Command-line license generator")
	fmt.Println()

	// Note
	fmt.Println("Note:")
	fmt.Println("  <license-name> refers to any license name listed in `license ls`")

	fmt.Println()

	// Usage
	fmt.Println("Usage:")
	for _, c := range []usageLine{
		{"help", "Show help information"},
		{"<license-name> [<args>]", "Generate the specified license"},
		{"ls", "List locally available licenses"},
		{"ls-remote", "List remote licenses"},
		{"update [--verbose]", "Update local licenses to latest remote versions"},
		{"version", "Print the current version"},
	} {
		fmt.Println(&c)
	}

	fmt.Println()

	// Examples
	fmt.Println("Example:")
	for _, c := range []exampleLine{
		{"license mit", "Print MIT license"},
		{"license mit --year 2050 --name Alice", "Print MIT license and override year and name"},
		{"license isc -o LICENSE.txt", "Write ISC license to a file named LICENSE.txt"},
	} {
		fmt.Println(&c)
	}

	fmt.Println()

	return true
}

func main() {
	args := os.Args[1:]
	var success bool

	if len(args) < 1 {
		success = help()
	} else {
		first := args[0]

		switch first {
		// Help information
		case "-h":
			fallthrough
		case "--help":
			fallthrough
		case "help":
			success = help()

		// Version information
		case "-v":
			fallthrough
		case "--version":
			fallthrough
		case "version":
			success = version()

		// List local licenses
		case "ls":
			fallthrough
		case "list":
			success = listLocal()

		// List remote licenses
		case "ls-remote":
			fallthrough
		case "list-remote":
			success = listRemote()

		// Update to latest remote licenses
		case "update":
			success = update()

		// Generate license if arg is known license name
		default:
			success = generate(args)
		}
	}

	if !success {
		os.Exit(1)
	}

	os.Exit(0)
}
