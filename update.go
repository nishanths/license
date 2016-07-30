package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/pkg/license"
)

func update() {
	err := doUpdate()
	if err != nil {
		errLogger.Println("failed to update licenses:", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func prepareTempDir() (string, error) {
	tempRoot, err := ioutil.TempDir("", "license-")
	if err != nil {
		return "", err
	}

	data := path.Join(tempRoot, "data")
	tmpl := path.Join(data, "tmpl")

	if err := os.MkdirAll(tmpl, 0700); err != nil {
		os.RemoveAll(tempRoot)
		return "", err
	}

	return tempRoot, nil
}

func doUpdate() error {
	// Exit early if we cannot find home directory.

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Work in a temporary directory, then
	// move updated licenses to the home directory.

	tempRoot, err := prepareTempDir()
	if err != nil {
		return nil
	}
	defer os.RemoveAll(tempRoot)

	c := license.NewClient()
	c.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	c.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

	// Write licenses.json.

	lics, err := c.List()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(tempRoot, "data", "licenses.json"))
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(lics); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	errs := make([]error, len(lics))

	for i, l := range lics {
		i, l := i, l
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Fetch full licenses.

			l, err := c.Info(l.Key)
			if err != nil {
				errs[i] = err
				return
			}

			// Write template file.

			content := textTemplate(l.Body)
			tmplFile, err := os.Create(filepath.Join(
				tempRoot, "data", "tmpl", strings.ToLower(l.Key)+".tmpl",
			))
			if err != nil {
				errs[i] = err
				return
			}
			if _, err := tmplFile.WriteString(content); err != nil {
				errs[i] = err
				return
			}
		}()
	}

	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	licensePath := path.Join(home, ".license")
	if err := os.RemoveAll(licensePath); err != nil {
		return err
	}
	return os.Rename(tempRoot, licensePath)
}

var (
	// placeholders is a map from placeholders in the JSON
	// to placeholders we use in templates.
	// The keys should be kept in sync with the regex below.
	placeholders = map[string]string{
		"year":     "Year",
		"fullname": "Name",
	}
	placeholdersRx = regexp.MustCompile("\\[(year|fullname)\\]")
)

func textTemplate(s string) string {
	return placeholdersRx.ReplaceAllStringFunc(s, func(m string) string {
		if s := placeholdersRx.FindStringSubmatch(m); len(s) > 0 {
			k := s[1]
			if v, ok := placeholders[k]; ok {
				return "{{." + v + "}}"
			}
		}
		return m
	})
}
