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
	"github.com/nishanths/license/remote"
	"io/ioutil"
	"os"
	"path"
	// "text/template"
)

func permissionsFailed(p string) {
	console.Error(fmt.Sprintf("Could not access %s. Make sure you have the permissions", p))
}

func created(p string) {
	console.Info(fmt.Sprintf("Created %s", p))
}

func generate(n string) {

}

func update() {

}

func listRemote() {
	a, _ := remote.List()
	fmt.Println(a)
}

func list() {
	a, _ := base.List()
	fmt.Println(a)
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

	// TODO: remove exisiting data only after successful bootstrap
	defer func() {
		fmt.Println("Bootstrap finished")
	}()

	licensePath := path.Join(home, base.LicenseDirectory)
	dataPath := path.Join(licensePath, base.DataDirectory)
	rawPath := path.Join(dataPath, base.RawDirectory)
	templatesPath := path.Join(dataPath, base.TemplatesDirectory)
	listFilePath := path.Join(dataPath, base.ListFile)

	if err := os.RemoveAll(licensePath); err != nil && os.IsPermission(err) {
		permissionsFailed(licensePath)
		return
	}

	// Create data directories
	pathsToMake := []string{rawPath, templatesPath}
	for _, p := range pathsToMake {
		if err := os.MkdirAll(p, 0700); err != nil {
			permissionsFailed(p)
			return
		}
		created(p)
	}

	// Fetch list from remote
	licenses, err := remote.List()
	if err != nil {
		console.Error("Failed to make licenses list from api.github.com")
		return
	}

	// Serialize the list into JSON
	// Write the serialized JSON to the list file
	serialized, err := json.Marshal(licenses)
	if err != nil {
		console.Error(fmt.Sprintf("Failed to serialize licenses. Please file an issue: %s", base.RepositoryIssuesURL))
		return
	}

	if err := ioutil.WriteFile(listFilePath, serialized, 0700); err != nil {
		console.Error(fmt.Sprintf("Failed to write %s", listFilePath))
		return
	}
	created(listFilePath)

	fullLicenses := make([]base.License, 0)

	// Fetch each license's full info
	// Serialize, then write individual files to disk
	for _, l := range licenses {
		full, err := remote.Info(&l)
		fmt.Fprintf(os.Stdout, full.Body)
		if err != nil {
			console.Error("Failed to fetch detailed license info")
			return
		}

		fullLicenses = append(fullLicenses, *full)

		serialized, err := json.Marshal(full)

		if err != nil {
			console.Error(fmt.Sprintf("Failed to serialize licenses. Please file an issue: %s", base.RepositoryIssuesURL))
			return
		}

		filePath := path.Join(rawPath, l.Key+".json")
		if err := ioutil.WriteFile(filePath, serialized, 0700); err != nil {
			console.Error(fmt.Sprintf("Failed to write %s", filePath))
			return
		}
	}

	fmt.Println(fullLicenses[0].TextTemplate())

	// if err := os.RemoveAll(home + LicenseDirectory); err != nil && os.IsExist(err) {
	// 	fmt.Println("did not find directory")
	// 	// TODO:
	// 	return
	// }

	// TODO: put files into data directory
}

func main() {
	bootstrap()
	// listRemote()
	// list()
	// tmpl, _ := template.New("sample").ParseFiles("sample")
	// tmpl.Execute(os.Stdout, nil)
}
