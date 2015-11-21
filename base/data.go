package base

import (
	"encoding/json"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

var Placeholders = map[string]string{
	"year":           "Year",
	"fullname":       "Name",
	"name of author": "Name",
}

var PlaceholdersRx *regexp.Regexp

func init() {
	keys := make([]string, 0)
	for key := range Placeholders {
		keys = append(keys, key)
	}
	PlaceholdersRx = regexp.MustCompile("[<{\\[]{1,}(" + strings.Join(keys, "|") + ")[>}\\]]{1,}")
}

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

func List() ([]License, error) {
	licenses := make([]License, 0)

	contents, err := read(path.Join(LicenseDirectory, DataDirectory, ListFile))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(contents, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}
