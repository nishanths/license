package base

import (
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"
)

// read returns the contents of a filename or path relative to the data directory.
func read(f string) ([]byte, error) {
	home, err := homedir.Dir()

	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(filepath.Join(home, LicenseDirectory, DataDirectory, f))

	if err != nil {
		return nil, err
	}

	return contents, nil
}

// readIndex reads the local index JSON file that has the list
// of current local licenses.
func readIndex() ([]byte, error) {
	return read(IndexFile)
}

// readFullInfo reads the local full JSON information for the given license.
func (l *License) readFullInfo() ([]byte, error) {
	return read(filepath.Join(RawDirectory, l.Key+".json"))
}

// readTemplate reads the template data and returns a template
// for a given license key.
func readTemplate(key string) (*template.Template, error) {
	home, err := homedir.Dir()

	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(path.Join(home, LicenseDirectory, DataDirectory, TemplatesDirectory, key+".tmpl"))

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
