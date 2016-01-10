package base

import (
	"github.com/tcnksm/go-gitconfig"
	"gopkg.in/nishanths/go-hgconfig.v1"
	"os"
)

const (
	// NameEnvVariable is the environment variable to lookup in the process
	// of determining the author's name to use on the license.
	NameEnvVariable = "LICENSE_FULL_NAME"
	defaultName     = ""
)

// getName attempts to implicitly guess the name to use
// on the license. The function looks for a name in the following order:
//
//   * env variable corresponding to `NameEnvVariable`, if it exists
//   * user.name from local git config
//   * user.name from global git config
//   * ui.username from hg config
//   * default (empty string)
//
func getName() string {
	if name := os.Getenv(NameEnvVariable); name != "" {
		return name
	} else if name, err := gitconfig.Username(); err == nil {
		return name
	} else if name, err := gitconfig.Global("user.name"); err == nil {
		return name
	} else if name, err := hgconfig.Username(); err == nil {
		return name
	} else {
		return defaultName
	}
}
