package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	appdir "github.com/ProtonMail/go-appdir"
	hgconfig "github.com/nishanths/go-hgconfig"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

const (
	nameEnv       = "LICENSE_FULL_NAME"
	versionString = "1.0.1"
	usageString   = "license [FLAGS] [LICENSE NAME]"
	helpString    = `usage: ` + usageString + `

Flags:
       -auth     GitHub credentials in format "username:token" for updating licenses (optional)
       -help     print help information
       -list     list available licenses
   -n, -name     full name on license (default %q)
   -o, -output   output filename (prints to stdout if not specified)
       -update   update all licenses to latest from GitHub
       -version  print version
   -y, -year     year on license (default %q)

Examples:
  license mit
  license -name Alice bsd-3-clause
  license -o LICENSE.txt mpl-2.0`
)

var (
	logger     = log.New(os.Stdout, "", 0)
	errLogger  = log.New(os.Stderr, "", log.Lshortfile)
	appDataDir = appdir.New("license").UserData()
	flags      = struct {
		License string // Type of license.
		Name    string // Name on license.
		Year    string // Year on license.
		Output  string // Output file.
		Auth    string
		Version bool
		Help    bool
		List    bool
		Update  bool
	}{}
)

func setupFlags() {
	flag.StringVar(&flags.Name, "name", name(), "name on license")
	flag.StringVar(&flags.Name, "n", name(), "name on license")

	flag.StringVar(&flags.Year, "year", year(), "year on license")
	flag.StringVar(&flags.Year, "y", year(), "year on license")

	flag.StringVar(&flags.Output, "output", "", "path to output file")
	flag.StringVar(&flags.Output, "o", "", "path to output file")

	flag.StringVar(&flags.Auth, "auth", "", "GitHub authentication")
	flag.BoolVar(&flags.Version, "version", false, "print version")
	flag.BoolVar(&flags.Help, "help", false, "print help")
	flag.BoolVar(&flags.List, "list", false, "print available licenses")
	flag.BoolVar(&flags.Update, "update", false, "get latest licenses")
}

func checkFlags() error {
	return nil
}

func main() {
	setupFlags()
	flag.Usage = usage
	flag.Parse()

	flags.License = strings.ToLower(flag.Arg(0))

	if flags.License == "" && !(flags.Update || flags.Version || flags.Help || flags.List) {
		help()
		os.Exit(1)
	}

	mainImpl()
}

func mainImpl() {
	// If any of the -version, -help, -list, or -update flags
	// is specified, the rest of the arguments is ignored.

	switch {
	case flags.Version:
		version()
	case flags.Help:
		help()
	case flags.List:
		list()
	case flags.Update:
		update()
	}

	generate()
	os.Exit(2)
}

func version() {
	errLogger.Printf("v%s\n", versionString)
	os.Exit(0)
}

func help() {
	errLogger.Printf(helpString+"\n", name(), year())
	os.Exit(0)
}

func usage() {
	errLogger.Println("usage: " + usageString)
	errLogger.Println("see 'license -help' for more details")
}

func name() string {
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

func year() string {
	return strconv.Itoa(time.Now().Year())
}

func ensureExists() error {
	_, err := os.Stat(appDataDir)
	if os.IsNotExist(err) {
		if e := doUpdate(); e != nil {
			return e
		}
		return nil
	}
	return err
}
