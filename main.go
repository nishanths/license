package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nishanths/license/base"
	"github.com/nishanths/license/local"
	"github.com/nishanths/license/remote"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/nishanths/simpleflag"
	shutil "github.com/termie/go-shutil"
)

const (
	applicationVersion = "0.1.0"
	helpIndent         = "  "
)

type helpLine struct {
	Command     string
	Description string
}

func (l *helpLine) String() string {
	return fmt.Sprintf("%s%-14s%s", helpIndent+helpIndent, l.Command, l.Description)
}

func failedToCreateDirectory(p string) {
	fmt.Println("license: failed to create directory", p)
}

func failedToCreateFile(p string) {
	fmt.Println("license: failed to create file", p)
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

func errorParsingArguments() {
	fmt.Println("license: error parsing arguments. See \"license help\".")
}

func unknownArgument(a string) {
	fmt.Printf("license: unknown argument: %v. See \"license help\".", a)
}

func badArgumentSyntax(a string) {
	fmt.Printf("license: bad argument: %v. See \"license help\"", a)
}

func generate(args []string) bool {
	if len(args) < 1 {
		fmt.Println("license: expected license name. See \"license help\".")
		return false
	}

	w := os.Stdout

	// Argument values
	var name, year, filename, licenseKey string

	nameCh := make(chan string, 1)
	go func(ch chan string) {
		ch <- base.GetName()
	}(nameCh)

	// Parse arguments
	generateFlagSet := simpleflag.NewFlagSet("generate")
	generateFlagSet.Add("name", []string{"--name", "-n"}, false)
	generateFlagSet.Add("year", []string{"--year", "-y"}, false)
	generateFlagSet.Add("output", []string{"--output", "-o"}, false)
	result, err := generateFlagSet.Parse(args)

	// Exit early if there is an error
	if err != nil {
		errorParsingArguments()
		return false
	}
	if len(result.BadFlags) > 0 {
		badArgumentSyntax(result.BadFlags[0])
		return false
	}
	if len(result.UnknownFlags) > 0 {
		for _, flag := range result.UnknownFlags {
			unknownArgument(flag)
		}
	}

	// Normalize
	// name
	if _, exists := result.Values["name"]; exists {
		name = result.Values["name"]
	} else {
		name = <-nameCh
	}

	// year
	if _, exists := result.Values["year"]; exists {
		year = result.Values["year"]
	} else {
		year = strconv.Itoa(time.Now().Year())
	}

	// filename
	filename = result.Values["output"]
	if filename != "" {
		var err error
		if w, err = os.Create(filename); err != nil {
			failedToCreateFile(filename)
			return false
		}
	}

	// Local licenses available currently
	licenses, err := local.List()
	if err != nil {
		// TODO: instead attempt bootstrap here or elsewhere
		listFailed()
		return false
	}

	// Get license key from remaining args
search:
	for _, arg := range result.Remaining {
		for _, license := range licenses {
			lowercasedArg := strings.ToLower(arg)
			if license.Key == lowercasedArg || license.Name == lowercasedArg {
				licenseKey = license.Key
				break search
			}
		}
	}

	if licenseKey == "" {
		fmt.Println("license: could not find license. See \"license ls\" for list of available licenses.")
		return false
	}

	tmpl, err := local.Template(licenseKey)
	if err != nil {
		fmt.Println("license: error loading template")
		fileAnIssue()
		return false
	}

	o := &base.Option{
		Name: name,
		Year: year,
	}

	// Execute
	if err := base.RenderTemplate(tmpl, o, w); err != nil {
		fmt.Println("license: error executing template")
		fileAnIssue()
		return false
	}

	return true
}

func help() bool {
	// Heading
	fmt.Println()
	fmt.Println(helpIndent + "Command-line license generator")
	fmt.Println()

	// Usage
	fmt.Println(helpIndent + "Usage:")
	fmt.Println()
	fmt.Println(helpIndent + helpIndent + "license [-y <year>] [-n <name>] [-o <filename>] [license-name]")
	fmt.Println()

	// Example
	fmt.Println(helpIndent + "Example:")
	fmt.Println()
	for _, c := range []helpLine{
		{"license mit", ""},
		{"license -y 2013 -n Alice mit", ""},
		{"license -o LICENSE.txt isc", ""},
	} {
		fmt.Println(&c)
	}
	fmt.Println()

	// Options
	fmt.Println(helpIndent + "Options:")
	fmt.Println()
	for _, c := range []helpLine{
		{"-y, --year", "Year to use on license"},
		{"-n, --name", "Name to use on license"},
		{"-o, --output", "Output file for license"},
	} {
		fmt.Println(&c)
	}
	fmt.Println()

	// Other commands:
	fmt.Println(helpIndent + "Other commands:")
	fmt.Println()
	for _, c := range []helpLine{
		{"ls", "List locally available license names"},
		{"ls-remote", "List remote license names"},
		{"update", "Update local licenses to latest remote versions"},
		{"help", "Show help information"},
		{"version", "Print current version"},
	} {
		fmt.Println(&c)
	}
	fmt.Println()

	// Note
	fmt.Println(helpIndent + "Run \"license ls\" for list of available license names")

	fmt.Println()

	return true
}

func listHelper(fn func() ([]base.License, error)) bool {
	licenses, err := fn()

	if err != nil {
		listFailed()
		return false
	}

	sort.Sort(base.ByLicenseKey(licenses))

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

func unknownCommand(args []string) bool {
	fmt.Printf("license: unknown command \"%s\". See \"license help\".\n", args[0])
	return true
}

func update(args []string) bool {
	// Determine if verbose logging is on
	verboseLog := false

	updateFlagSet := simpleflag.NewFlagSet("update")
	updateFlagSet.Add("verbose", []string{"--verbose", "-v"}, true)

	res, err := updateFlagSet.Parse(args)
	if err != nil {
		errorParsingArguments()
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
		case "--help":
			fallthrough
		case "help":
			success = help()

		// Version information
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

		default:
			success = generate(args)
		}
	}

	if !success {
		os.Exit(1)
	}

	os.Exit(0)
}
