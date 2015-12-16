package base

import (
	gitconfig "github.com/tcnksm/go-gitconfig"
	"os"
)

const (
	NameEnvVariable = "LICENSE_FULL_NAME"
	defaultName     = ""
)

func GetName() string {
	if name, exists := os.LookupEnv(NameEnvVariable); exists { // Use the env variable if it exists
		return name
	} else if name, err := gitconfig.Username(); err == nil { // Attempt to use local gitconfig for name
		return name
	} else if name, err := gitconfig.Global("user.name"); err == nil { // Attempt to use global gitconfig for name
		return name
	} else { // Finally use the default
		return defaultName
	}
}
