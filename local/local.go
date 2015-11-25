package local

import (
	"encoding/json"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"io/ioutil"
	"path"
	"sort"
	"text/template"
)

func read(filename string) ([]byte, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(path.Join(home, filename))
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func Template(key string) (*template.Template, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFiles(path.Join(home, base.LicenseDirectory, base.DataDirectory, base.TemplatesDirectory, key+".tmpl"))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func List() ([]base.License, error) {
	licenses := make([]base.License, 0)

	contents, err := read(path.Join(base.LicenseDirectory, base.DataDirectory, base.ListFile))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(contents, &licenses); err != nil {
		return nil, err
	}

	sort.Sort(base.ByLicenseKey(licenses))
	return licenses, nil
}

func Info(l *base.License) (base.License, error) {
	var r base.License

	contents, err := read(path.Join(base.LicenseDirectory, base.DataDirectory, base.RawDirectory, l.Key+".json"))
	if err != nil {
		return r, err
	}

	if err := json.Unmarshal(contents, &r); err != nil {
		return r, err
	}

	return r, nil
}
