package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type unknownLicenseError struct {
	License string
}

func (e unknownLicenseError) Error() string {
	return fmt.Sprintf(`error: unknown license %q
see 'license -list' for the list of available licenses`, e.License)
}

func generate() {
	err := doGenerate()
	if err != nil {
		errLogger.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func doGenerate() error {
	if err := ensureExists(); err != nil {
		return err
	}
	p := filepath.Join(appDataDir, "tmpl", flags.License+".tmpl")
	t, err := template.ParseFiles(p)
	if err != nil {
		if os.IsNotExist(err) {
			return unknownLicenseError{flags.License}
		}
		return err
	}

	var out io.Writer = os.Stdout
	if flags.Output != "" {
		f, err := os.Create(filepath.Clean(flags.Output))
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	// For non-nil error, data may have been written, but at least
	// we can set a non-zero exit code.
	return t.Execute(out, struct {
		Name string
		Year string
	}{flags.Name, flags.Year})
}
