// license generates the text for a license of your choice
// Usage: license <license-name>
// Example: license mit

package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"github.com/nishanths/license/console"
	"github.com/nishanths/license/remote"
	"os"
	"path"
	// "text/template"
)

func permissionsFailed(p string) {
	console.Error(fmt.Sprintf("Could not access %s. Make sure you have the permissions.", p))
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
	home, err := homedir.Dir()
	if err != nil {
		console.Error("Unable to locate home directory.")
		return
	}

	licensePath := path.Join(home, base.LicenseDirectory)
	if err := os.RemoveAll(licensePath); err != nil && os.IsPermission(err) {
		permissionsFailed(licensePath)
		return
	}
	console.Info(fmt.Sprintf("Removed %s", licensePath))

	pathsToMake := []string{path.Join(home, base.LicenseDirectory, base.DataDirectory, base.TemplatesDirectory)}
	for _, p := range pathsToMake {
		if err := os.MkdirAll(p, 0700); err != nil {
			permissionsFailed(p)
			return
		}
		console.Info(fmt.Sprintf("Created %s", p))
	}

	licenses, err := remote.List()
	if err != nil {
		console.Error("Failed to fetch licenses list from api.github.com.")
	}

	info, _ := remote.Info(&licenses[0])
	fmt.Println(info)

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
