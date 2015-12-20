package base

import (
	"gopkg.in/nishanths/simpleflag.v1"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type renderOption struct {
	Year string
	Name string
}

func renderTemplate(t *template.Template, o *renderOption, w io.Writer) error {
	return t.ExecuteTemplate(w, t.Name(), o)
}

// Generate parses arguments and outputs the selected license.
// Generate returns a non-nil error if it is unable to do so successfully.
func Generate(args []string) error {
	if len(args) < 1 {
		return NewErrExpectedLicenseName()
	}

	// use stdout as default writer for now, until
	// specfied in the args
	w := os.Stdout

	// arguments values
	var name, year, filename, licenseKey string

	// start looking for the default name
	// to use on the license, in case we need it
	nameCh := make(chan string, 1)
	go func(ch chan string) {
		ch <- getName()
	}(nameCh)

	defer func() {
		close(nameCh)
	}()

	// parse arguments
	generateFlagSet := simpleflag.NewFlagSet("generate")
	generateFlagSet.Add("name", []string{"--name", "-name", "-n"}, false)
	generateFlagSet.Add("year", []string{"--year", "-year", "-y"}, false)
	generateFlagSet.Add("output", []string{"--output", "-output", "-o"}, false)
	result, err := generateFlagSet.Parse(args)

	// exit early if there is an error
	// with args
	if err != nil {
		return NewErrParsingArguments()
	}

	if len(result.BadFlags) > 0 {
		return NewErrBadFlagSyntax(result.BadFlags[0])
	}

	// normalize:

	// 1. name
	if n, exists := result.Values["name"]; exists {
		name = n
	} else {
		name = <-nameCh
	}

	// 2. year
	if y, exists := result.Values["year"]; exists {
		year = y
	} else {
		year = strconv.Itoa(time.Now().Year())
	}

	// 3. filename
	filename = result.Values["output"]

	// get locally available licenses
	licenses, err := getLocalList()
	if err != nil {
		return NewErrReadFailed()
	}

	// find license key from remaining args
search:
	for _, arg := range result.Remaining {
		for _, license := range licenses {
			lowercasedArg := strings.ToLower(arg)
			if strings.ToLower(license.Key) == lowercasedArg || strings.ToLower(license.Name) == lowercasedArg {
				licenseKey = license.Key
				break search
			}
		}
	}

	if licenseKey == "" {
		return NewErrCannotFindLicense()
	}

	tmpl, err := readTemplate(licenseKey)

	if err != nil {
		return NewErrLoadingTemplate(licenseKey + ".tmpl")
	}

	o := &renderOption{
		Name: name,
		Year: year,
	}

	// create the file since we are close to succeeding
	if filename != "" {
		var err error
		if w, err = os.Create(filename); err != nil {
			return NewErrWriteFileFailed(filename)
		}
	}

	// execute template on file
	if err := renderTemplate(tmpl, o, w); err != nil {
		os.Remove(filename)
		return NewErrExecutingTemplate(tmpl)
	}

	return nil
}
