package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/nishanths/go-hgconfig"
	"github.com/tcnksm/go-gitconfig"
)

const (
	nameEnv       = "LICENSE_FULL_NAME"
	versionString = "2.0.0"

	usageString = `Usage: license [flags] [license-type]

Flags:
       -help     print help information
       -list     print list of available license types
   -n, -name     full name to use on license (default %q)
   -o, -output   path to output file (prints to stdout if unspecified)
       -version  print version
   -y, -year     year to use on license (default %q)

Examples:
  license mit 
  license -name "Alice L" -year 2013 bsd-3-clause
  license -o LICENSE.txt mpl-2.0`
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

var (
	fName    string
	fYear    string
	fOutput  string
	fVersion bool
	fHelp    bool
	fList    bool
)

func main() {
	flag.StringVar(&fName, "name", "", "name on license")
	flag.StringVar(&fName, "n", "", "name on license")
	flag.StringVar(&fYear, "year", "", "year on license")
	flag.StringVar(&fYear, "y", "", "year on license")
	flag.StringVar(&fOutput, "output", "", "path to output file")
	flag.StringVar(&fOutput, "o", "", "path to output file")
	flag.BoolVar(&fVersion, "version", false, "print version")
	flag.BoolVar(&fHelp, "help", false, "print help")
	flag.BoolVar(&fList, "list", false, "print available licenses")

	flag.Usage = func() {
		printUsage()
		os.Exit(1)
	}
	flag.Parse()

	run()
}

func run() {
	if flag.NArg() != 1 && !(fVersion || fHelp || fList) {
		printUsage()
		os.Exit(1)
	}

	switch {
	case fVersion:
		printVersion()
		os.Exit(0)

	case fHelp:
		printUsage()
		os.Exit(0)

	case fList:
		printList()
		os.Exit(0)

	default:
		license := strings.ToLower(flag.Arg(0))
		printLicense(license, fOutput, getName(), getYear()) // internally calls os.Exit() on failure
	}
}

func printVersion() {
	stdout.Printf("%s", versionString)
}

func printUsage() {
	stderr.Printf(usageString, getName(), getYear())
}

func getName() string {
	if fName != "" {
		return fName
	}
	n := os.Getenv(nameEnv)
	if n != "" {
		return n
	}
	n, err := gitconfig.Username()
	if err == nil {
		return n
	}
	n, err = gitconfig.Global("user.name")
	if err == nil {
		return n
	}
	n, err = hgconfig.Username()
	if err == nil {
		return n
	}
	usr, err := user.Current()
	if err == nil {
		return usr.Name
	}
	return ""
}

func getYear() string {
	if fYear != "" {
		return fYear
	}
	return strconv.Itoa(time.Now().Year())
}
