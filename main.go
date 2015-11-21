// license generates the text for a license of your choice
// Usage: license <license-name>
// Example: license mit

package main

import (
	"encoding/json"
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
	// "text/template"
)

const (
	tempDirectoryPrefix = "license-"
)

func permissionsFailed(p string) {
	console.Error(fmt.Sprintf("Could not access %s. Make sure you have the permissions", p))
}

func createPathFail(p string) {
	console.Error(fmt.Sprintf("Failed to make %s. Make sure you have the permissions", p))
}

func createPathSuccess(p string) {
	console.Info(fmt.Sprintf("Created %s", p))
}

func generate(n string) {

}

func update() {

}

func list(fn func() ([]base.License, error)) {
	licenses, err := fn()
	if err != nil {
		console.Error(fmt.Sprintf("Failed to list licenses. Run `$ license bootstrap` before trying again. Otherwise, please create an Issue: %s", base.RepositoryIssuesURL))
		return
	}
	base.RenderList(&licenses)
}

func listRemote() {
	list(remote.List)
}

func listLocal() {
	list(local.List)
}

func help() {

}

func bootstrap() {
	// Bail immediately if we cannot find the user's home directory
	home, err := homedir.Dir()
	if err != nil {
		console.Error("Unable to locate home directory")
		return
	}

	tempLicensePath, err := ioutil.TempDir(os.TempDir(), tempDirectoryPrefix)
	if err != nil {
		createPathFail("temporary directory")
		return
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
			return
		}
	}
	createPathSuccess(fmt.Sprintf("temp directory %s", tempLicensePath))

	// Fetch list from remote
	licenses, err := remote.List()
	if err != nil {
		console.Error("Failed to make licenses list from api.github.com")
		return
	}
	console.Info("Fetched license list from api.github.com")

	// Serialize the list into JSON
	// Write the serialized JSON to the list file
	serialized, err := json.Marshal(licenses)
	if err != nil {
		console.Error(fmt.Sprintf("Failed to serialize licenses. Please create an Issue: %s", base.RepositoryIssuesURL))
		return
	}

	if err := ioutil.WriteFile(listFilePath, serialized, 0700); err != nil {
		createPathFail(listFilePath)
		return
	}

	fullLicenses := make([]base.License, 0)

	// Fetch each license's full info
	// - Serialize to JSON and write to disk
	// - Convert to text template and write to disk
	for _, l := range licenses {
		fullLicense, err := remote.Info(&l)
		if err != nil {
			console.Error("Failed to make detailed license info")
			return
		}

		fullLicenses = append(fullLicenses, fullLicense)

		serialized, err := json.Marshal(fullLicense)
		if err != nil {
			console.Error(fmt.Sprintf("Failed to serialize licenses. Please file an issue: %s", base.RepositoryIssuesURL))
			return
		}

		rawFilePath := path.Join(rawPath, l.Key+".json")
		if err := ioutil.WriteFile(rawFilePath, serialized, 0700); err != nil {
			createPathFail(rawFilePath)
			return
		}

		templateData := fullLicense.TextTemplate()
		templateFilePath := path.Join(templatesPath, l.Key+".tmpl")
		if err := ioutil.WriteFile(templateFilePath, []byte(templateData), 0700); err != nil {
			createPathFail(templateFilePath)
			return
		}
	}

	// Remove exisitng path
	realLicensePath := path.Join(home, base.LicenseDirectory)
	if err := os.RemoveAll(realLicensePath); err != nil && os.IsPermission(err) {
		permissionsFailed(realLicensePath)
		return
	}

	// Copy temp data to real path
	if err := shutil.CopyTree(tempLicensePath, realLicensePath, nil); err != nil {
		console.Error(fmt.Sprintf("Failed to copy data to %s", realLicensePath))
		return
	}
	createPathSuccess(fmt.Sprintf("and copied data to %s", realLicensePath))
}

func render() {
	var c base.Config
	c.Prepare("Nishanth Shanmugham", "")
	o := base.NewOption(c.Name)

	var l base.License
	l.Key = "isc"

	tmpl, _ := local.Template(&l)
	base.RenderTemplate(tmpl, &o)
}

func main() {
	// listRemote()
	// listLocal()
	render()
}
