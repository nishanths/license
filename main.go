package main

import (
	"encoding/json"
	"flag"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"github.com/nishanths/license/local"
	"github.com/nishanths/license/remote"
	"github.com/nishanths/simpleflag"
	shutil "github.com/termie/go-shutil"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	applicationVersion = "0.1.0"
	indent             = "  "
)

type usageLine struct {
	Command     string
	Description string
}

type exampleLine usageLine

func (l *usageLine) String() string {
	return fmt.Sprintf("%s%-12s%s", indent+indent, l.Command, l.Description)
}

func (l *exampleLine) String() string {
	return fmt.Sprintf("%s%-40s", indent+indent, l.Command)
}

func failedToCreateDirectory(p string) {
	fmt.Println("license: failed to create directory", p)
}

func failedToSerializeLicenses() {
	fmt.Println("license: failed to serialize licenses")
}

func fileAnIssue() {
	fmt.Println("license: please create an issue at", base.RepositoryIssuesURL)
}

func listFailed() {
	fmt.Println("license: failed to make list")
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
		f, err = os.Create(output)
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

func help() bool {
	// Heading
	fmt.Println()
	fmt.Println(indent + "Command-line license generator")
	fmt.Println()

	// Usage
	fmt.Println(indent + "Usage:")
	fmt.Println()
	for _, c := range []usageLine{
		{"use", "Create LICENSE file with specified license"},
		{"view", "Print specified license on stdout"},
		{"ls", "List locally available license names"},
		{"ls-remote", "List remote license names"},
		{"update", "Update local licenses to latest remote versions"},
		{"help", "Show help information"},
		{"version", "Print current version"},
	} {
		fmt.Println(&c)
	}
	fmt.Println()

	// Example
	fmt.Println(indent + "Example:")
	fmt.Println()
	for _, c := range []exampleLine{
		{"license use mit", "Create MIT license in file named `LICENSE`"},
		{"license use mit --filename LICENSE.md", "Create MIT license in file named `LICENSE.md`"},
		{"license use isc --year 2050 --name Alice", "Use custom year and name"},
	} {
		fmt.Println(&c)
	}
	fmt.Println()

	// Note
	fmt.Println(indent + "Run `license ls` to see list of available licenses")
	fmt.Println()

	return true
}

func listHelper(fn func() ([]base.License, error)) bool {
	licenses, err := fn()
	if err != nil {
		listFailed()
		return false
	}

	fmt.Println()
	fmt.Println("  Available licenses:\n")
	base.RenderList(&licenses)
	fmt.Println()
	return true
}

func listLocal() bool {
	return listHelper(local.List)
}

func listRemote() bool {
	return listHelper(remote.List)
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

func unknownCommand(args []string) bool {
	fmt.Printf("license: unknown command `%s`. See `license help`.\n", args[0])
	return true
}

func update(args []string) bool {
	// Determine if verbose logging is on
	verboseLog := false

	updateFlagSet := simpleflag.NewFlagSet("update")
	updateFlagSet.Add("verbose", []string{"--verbose", "-v"}, false, "")
	res, err := updateFlagSet.Parse(args)
	if err != nil {
		fmt.Println("license: error parsing arguments")
		return false
	}

	if _, exists := res.Values["verbose"]; exists {
		verboseLog = true
	}

	// Bail immediately if we cannot find the user's home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("license: unable to locate home directory")
		return false
	}

	// Create temporary directory
	tempLicensePath, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		fmt.Println("license: failed to create temporary directory")
		return false
	}

	dataPath := path.Join(tempLicensePath, base.DataDirectory)
	rawPath := path.Join(dataPath, base.RawDirectory)
	templatesPath := path.Join(dataPath, base.TemplatesDirectory)
	listFilePath := path.Join(dataPath, base.ListFile)

	defer func() {
		if verboseLog {
			fmt.Println("cleaning up")
			fmt.Println("removing temporary directory", tempLicensePath)
		}
		os.RemoveAll(tempLicensePath)
		if verboseLog {
			fmt.Println("bootstrap complete!")
		}
	}()

	// Create data directories
	pathsToMake := []string{rawPath, templatesPath}
	for _, p := range pathsToMake {
		if err := os.MkdirAll(p, 0700); err != nil {
			failedToCreateDirectory(p)
			return false
		}
	}

	// Fetch list from remote
	licenses, err := remote.List()
	if err != nil {
		fmt.Println("license: failed to get license list from api.github.com")
		return false
	}

	if verboseLog {
		fmt.Println("fetched license list from api.github.com")
	}

	// Serialize the list into JSON
	// Write the serialized JSON to the list file
	serialized, err := json.MarshalIndent(licenses, "", "    ")
	if err != nil {
		failedToSerializeLicenses()
		fileAnIssue()
		return false
	}

	if err := ioutil.WriteFile(listFilePath, serialized, 0700); err != nil {
		failedToCreateDirectory(listFilePath)
		return false
	}

	fullLicenses := make([]base.License, 0)

	// Fetch each license's full info
	// - Serialize to JSON and write to disk
	// - Convert to text template and write to disk
	for _, l := range licenses {
		fullLicense, err := remote.Info(&l)
		if err != nil {
			fmt.Println("license: failed to make detailed license info")
			return false
		}

		fullLicenses = append(fullLicenses, fullLicense)

		serialized, err := json.MarshalIndent(fullLicense, "", "    ")
		if err != nil {
			failedToSerializeLicenses()
			fileAnIssue()
			return false
		}

		rawFilePath := path.Join(rawPath, l.Key+".json")
		if err := ioutil.WriteFile(rawFilePath, serialized, 0700); err != nil {
			failedToCreateDirectory(rawFilePath)
			return false
		}

		templateData := fullLicense.TextTemplate()
		templateFilePath := path.Join(templatesPath, l.Key+".tmpl")
		if err := ioutil.WriteFile(templateFilePath, []byte(templateData), 0700); err != nil {
			failedToCreateDirectory(templateFilePath)
			return false
		}
	}

	// Remove exisitng path
	realLicensePath := path.Join(home, base.LicenseDirectory)
	if err := os.RemoveAll(realLicensePath); err != nil && os.IsPermission(err) {
		fmt.Println("license: failed to access directory", realLicensePath)
		fmt.Println("license: make sure you have the right permissions")
		return false
	}

	// Copy temp data to real path
	if err := shutil.CopyTree(tempLicensePath, realLicensePath, nil); err != nil {
		fmt.Println("license: failed to copy data to", realLicensePath)
		return false
	}

	if verboseLog {
		fmt.Println("copied data to", realLicensePath)
	}

	return true
}

func version() bool {
	fmt.Println(applicationVersion)
	return true
}

func main() {
	args := os.Args[1:]
	var success bool

	if len(args) < 1 {
		success = help()
	} else {
		command := args[0]

		switch command {
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
			fallthrough
		case "bootstrap":
			success = update(args[1:])

		// Generate license
		case "view":
			fallthrough
		case "use":
			success = generate(args[1:])

		// Show help if unknown command
		default:
			success = unknownCommand(args)
		}
	}

	if !success {
		os.Exit(1)
	}

	os.Exit(0)
}
