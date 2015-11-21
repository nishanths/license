package base

import (
	"encoding/json"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"path"
)

func List() ([]License, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(path.Join(home, LicenseDirectory, DataDirectory, ListFile))
	if err != nil {
		return nil, err
	}

	licenses := make([]License, 0)

	if err := json.Unmarshal(contents, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}
