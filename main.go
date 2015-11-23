// license generates the text for a license of your choice
// Usage: license <license-name>
// Example: license mit

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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	// "text/template"
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

func generate(target string) bool {
	licenses, err := local.List()
	if err != nil {
		listFailed()
		return false
	}

	for _, l := range licenses {
		if l.Key == target || l.Name == target {
			return render(l.Key)
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
	fmt.Println("v" + applicationVersion)
	return true
}

func list(args []string) bool {
	listSet := flag.NewFlagSet("list", flag.ExitOnError)
	var local bool
	listSet.BoolVar(&local, "local", true, "List all licenses stored locally")
	remote := listSet.Bool("remote", false, "List all licenses from api.github.com")

	err := listSet.Parse(args)
	if err != nil {
		return false
	}

	if *remote {
		return listRemote()
	} else if local {
		return listLocal()
	}

	return false
}

func listHelper(fn func() ([]base.License, error), ms time.Duration) bool {
	licenses, err := fn()
	if err != nil {
		listFailed()
		return false
	}

	base.RenderList(&licenses, ms)
	return true
}

func listRemote() bool {
	return listHelper(remote.List, 0)
}

func listLocal() bool {
	return listHelper(local.List, 0)
}

type helpline struct {
	Command     string
	Description string
}

func (l *helpline) String() string {
	return fmt.Sprintf("   %-12s%-10s", l.Command, l.Description)
}

func help() bool {
	commands := []helpline{
		{"make", "Generate the specified license"},
		{"config", "Configure application globals"},
		{"help", "Show help information"},
		{"list", "List available licenses"},
		{"update", "Update local data with latest online licenses"},
		{"version", "Print current version"},
	}

	fmt.Printf("usage: license <license-name|command> [<args>]\n\n")
	fmt.Printf("Available commands:\n")
	for _, c := range commands {
		fmt.Println(&c)
	}
	fmt.Printf("\nSee 'license help <command>' to learn about a particular command\n")

	return true
}

func config(args []string) bool {
	configSet := flag.NewFlagSet("config", flag.ExitOnError)
	name := configSet.String("name", "", "Set full name on future licenses") // no default name

	if len(args) < 1 {
		console.Error("license: expected options for config")
		configSet.PrintDefaults()
		return false
	}

	configSet.Parse(args)

	c := base.Config{Name: *name}
	data, err := json.MarshalIndent(c, "", "    ")

	if err != nil {
		console.Error("license: failed to serialize config")
		return false
	}

	home, err := homedir.Dir()
	if err != nil {
		console.Error("license: unable to locate home directory")
		return false
	}

	file := path.Join(home, base.LicenseConfigFile)
	if err := ioutil.WriteFile(file, data, 0700); err != nil {
		console.Error(fmt.Sprintf("license: failed to create config file: %s", file))
		return false
	}

	return true
}

func render(key string) bool {
	var c base.Config
	c.Prepare("", "")
	o := base.NewOption(c.Name)

	tmpl, err := local.Template(key)
	if err != nil {
		// TODO: error message
		return false
	}

	base.RenderTemplate(tmpl, &o)
	return true
}

func main() {
	args := os.Args[1:]

	var success bool

	if len(args) < 1 {
		success = help()
	} else {
		first := args[0]

		switch strings.ToLower(first) {
		case "config":
			success = config(args[1:])
		case "help":
			success = help()
		case "version":
			success = version()
		case "list":
			success = list(args[1:])
		case "update":
			success = update()
		case "use":
			fallthrough
		case "generate":
			fallthrough
		case "make":
			if len(args) < 2 {
				fmt.Println("license: expected: license name")
			} else {
				licenseName := strings.ToLower(strings.TrimSpace(args[1]))
				if generate(licenseName) {
					os.Exit(0)
				} else {
					os.Exit(1)
				}
			}
		}
	}

	if !success {
		os.Exit(1)
	}

	os.Exit(0)

	// You can get individual args with normal indexing.
	// arg := os.Args[3]

	fmt.Println(args)
	// fmt.Println(arg)

	name := flag.String("name", "", "Full name on license")
	year := flag.Int("year", time.Now().Year(), "Year on license")

	// flag.NewFlagSet("name", flag.ExitOnError)
	flag.Parse()

	fmt.Println("name:", *name)
	fmt.Println("year:", *year)
	fmt.Println("tail:", flag.Args())
}
