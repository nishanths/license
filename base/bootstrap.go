package base

import (
	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/logger"
	"github.com/termie/go-shutil"
	"gopkg.in/nishanths/simpleflag.v1"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func setLogLevel(args []string) error {
	flagSet := simpleflag.NewFlagSet("")
	flagSet.Add("quiet", []string{"--quiet", "-quiet", "-q"}, true)
	flagSet.Add("verbose", []string{"--verbose", "-verbose", "-v"}, true)
	result, err := flagSet.Parse(args)

	if err != nil {
		return NewErrParsingArguments()
	}

	if len(result.BadFlags) > 0 {
		return NewErrBadFlagSyntax(result.BadFlags[0])
	}

	if _, exists := result.Values["quiet"]; exists {
		logger.SetQuiet(true)
	}

	if _, exists := result.Values["verbose"]; exists {
		logger.SetVerbose(true)
	}

	return nil
}

// Bootstrap updates local licenses
// to the latest online versions
func Bootstrap(args []string) error {
	if err := setLogLevel(args); err != nil {
		return err
	}

	// bail immediately if we cannot find the user's home directory
	home, err := homedir.Dir()
	if err != nil {
		return NewErrCannotLocateHomeDir()
	}

	// create temporary directory
	tempLicensePath, err := ioutil.TempDir("", tempDirPrefix)
	if err != nil {
		return NewErrCreateTempDirFailed(tempLicensePath)
	}

	// make path strings relative to temp directory
	dataPath := path.Join(tempLicensePath, DataDirectory)
	rawPath := path.Join(dataPath, RawDirectory)
	templatesPath := path.Join(dataPath, TemplatesDirectory)
	indexFilePath := filepath.Join(dataPath, IndexFile)

	// defer cleaning up of temporary directory
	defer func() {
		os.RemoveAll(tempLicensePath)
	}()

	// create data directories
	pathsToMake := []string{rawPath, templatesPath}

	for _, p := range pathsToMake {
		if err := os.MkdirAll(p, perm); err != nil {
			return NewErrCreateDirFailed(p)
		}
	}

	// fetch index file json
	// return error if we failed to fetch
	serialized, err := fetchIndex()
	if err != nil {
		return NewErrFetchFailed()
	}

	logger.VerbosePrintln("fetched data from api.github.com...")

	// write fetched index JSON to file
	if err := ioutil.WriteFile(indexFilePath, serialized, perm); err != nil {
		return NewErrCreateDirFailed(indexFilePath)
	}

	logger.VerbosePrintln("created local index file...")

	// make list of short licenses
	// from the fetched index file
	licenses, err := jsonToList(serialized)

	if err != nil {
		return NewErrDeserializeFailed(serialized)
	}

	// TODO: concurrent fetch
	for _, l := range licenses {
		// fetch full license info JSON
		content, err := l.fetchFullInfo()
		if err != nil {
			return NewErrFetchFailed()
		}

		// write JSON to disk
		rawFilePath := filepath.Join(rawPath, l.Key+".json")
		if err := ioutil.WriteFile(rawFilePath, serialized, perm); err != nil {
			return NewErrWriteFileFailed(rawFilePath)
		}

		// deserialize JSON to License struct
		fullLicense, err := jsonToLicense(content)
		if err != nil {
			return NewErrDeserializeFailed(content)
		}

		// construct template and save template in templates directory
		templateData := textTemplateString(&fullLicense)
		templateFilePath := filepath.Join(templatesPath, l.Key+".tmpl")

		if err := ioutil.WriteFile(templateFilePath, []byte(templateData), perm); err != nil {
			return NewErrWriteFileFailed(templateFilePath)
		}
	}

	logger.VerbosePrintln("created license templates...")

	// remove exisiting path + data
	realLicensePath := path.Join(home, LicenseDirectory)

	if err := os.RemoveAll(realLicensePath); err != nil && os.IsPermission(err) {
		return NewErrRemovePathFailed(realLicensePath)
	}

	// copy temp data to real path
	if err := shutil.CopyTree(tempLicensePath, realLicensePath, nil); err != nil {
		return NewErrCopyTreeFailed(tempLicensePath, realLicensePath)
	}

	logger.VerbosePrintln("bootstrap complete!")

	return nil
}
