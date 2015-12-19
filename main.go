package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/base"
	"os"
	"path"
	"sync"
	"time"
)

// pathExists returns true if the path exists.
func pathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

// main returns exit code 0 on success
// and exit code 1 on error.
func main() {
	args := os.Args[1:]
	var wg sync.WaitGroup
	var mainErr error

	// * Check existence of license data directory
	// and start making it if it is not present.
	// * Also, another time we update is around once every 20 runs
	// so that the licenses list is up to date.
	// * If we cannot find the home directory, we silently ignore the issue
	// for now; the specific command function will return an error when called.
	if home, err := homedir.Dir(); err == nil {
		updateRequired := (time.Now().Unix() % 20) == 0
		bootstrapRequired := !pathExists(path.Join(home, base.LicenseDirectory, base.DataDirectory))
		repetitiveCommand := len(args) >= 1 && !(args[0] == "update" || args[0] == "bootstrap")

		if (updateRequired || bootstrapRequired) && !(repetitiveCommand) {
			wg.Add(1)

			go func() {
				defer wg.Done()
				base.Bootstrap([]string{"--quiet"})
			}()
		}
	}

	if len(args) < 1 {
		mainErr = base.Help()
	} else {
		command := args[0]

		switch command {
		// Help information
		case "--help":
			fallthrough
		case "help":
			mainErr = base.Help()

		// Version information
		case "--version":
			fallthrough
		case "version":
			mainErr = base.Version()

		// Update to latest remote licenses
		case "update":
			fallthrough
		case "bootstrap":
			mainErr = base.Bootstrap(args[1:])

		// List remote licenses
		case "ls-remote":
			fallthrough
		case "list-remote":
			mainErr = base.ListRemote()

		// List local licenses
		case "ls":
			fallthrough
		case "list":
			wg.Wait()
			mainErr = base.ListLocal()

		default:
			wg.Wait()
			mainErr = base.Generate(args)
		}
	}

	wg.Wait()

	if mainErr != nil {
		fmt.Fprintln(os.Stderr, mainErr)
		os.Exit(1)
	}

	os.Exit(0)
}
