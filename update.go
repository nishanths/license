package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/pkg/license"
	shutil "github.com/termie/go-shutil"
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

	data := filepath.Join(tempRoot, "data")
	tmpl := filepath.Join(data, "tmpl")

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
	// copy to the home directory.

	tempRoot, err := prepareTempDir()
	if err != nil {
		return nil
	}
	defer os.RemoveAll(tempRoot)

	c := license.NewClient()
	if flags.Auth != "" {
		p := strings.Split(flags.Auth, ":")
		if len(p) == 2 {
			c.Username = p[0]
			c.Token = p[1]
		} else {
			return errors.New(`license: auth should be in format "username:token"`)
		}
	} else {
		c.ClientID = os.Getenv("GITHUB_CLIENT_ID")
		c.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	}

	// Write licenses.json.

	lics, err := c.List()
	if err != nil {
		return handleAPIError(err)
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
				errs[i] = handleAPIError(err)
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

	licensePath := filepath.Join(home, ".license")
	if err := os.RemoveAll(licensePath); err != nil {
		return err
	}
	return shutil.CopyTree(tempRoot, licensePath, nil)
}

// handleAPIError handles HTTP errors as a special case:
//   - simply logging the error may print the client ID
//     and client secret as part of the error message URL.
//   - if err is StatusError (403, API rate limit...)
//     then print -auth flag instructions.
func handleAPIError(err error) error {
	switch e := err.(type) {
	case nil:
		return nil
	case license.StatusError:
		if e.StatusCode == 403 && strings.Contains(e.Details.Message, "API rate limit") {
			return fmt.Errorf("%s\n\n%s\n%s", e.Error(),
				`Use a GitHub personal access token with the "-auth" flag.`,
				`See https://github.com/settings/tokens and "license -help" for more details.`)
		}
		return err
	default:
		return errors.New("error: failed to fetch licenses")
	}
}

var (
	// placeholders is a map from placeholders in GitHub's JSON
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
