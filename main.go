// license generates the text for a license of your choice
// Usage: license <license-name>
// Example: license mit

package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	// "github.com/nishanths/license/base"
	"github.com/nishanths/license/remote"
	"os"
	// "text/template"
)

const (
	LICENSE_CONFIG_FILE = ".licenseconfig"
	LICENSE_DIRECTORY   = ".license"
)

func generate() {

}

func update() {

}

func listRemote() {
	ch := make(chan *remote.LicensesResult)
	go remote.FetchLicensesList(ch)

	res := <-ch

	if res.Error != nil {
		fmt.Println("ls-remote had a fetch result error")
		return
	}

	// TODO: format
	fmt.Println(res.Licenses)
}

func list() {

}

func help() {

}

func bootstrap() {
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println("did not find home")
		// TODO: bootstrap failed
		return
	}

	err = os.Remove(home + LICENSE_CONFIG_FILE)
	if err != nil && os.IsExist(err) {
		fmt.Println("did not find config dotfile")
		// TODO:
		return
	}

	err = os.RemoveAll(home + LICENSE_DIRECTORY)
	if err != nil && os.IsExist(err) {
		fmt.Println("did not find directory")
		// TODO:
		return
	}

	// TODO
}

func main() {
	bootstrap()
	listRemote()
	// tmpl, _ := template.New("sample").ParseFiles("sample")
	// tmpl.Execute(os.Stdout, nil)
}
