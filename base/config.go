package base

import (
	// "encoding/json"
	"fmt"
	gitconfig "github.com/tcnksm/go-gitconfig"
	"os"
)

const (
	FullnameEnvVariable = "LICENSE_FULLNAME"
	defaultName         = ""
)

type Config struct {
	Name string `json:"name"`
}

func NewConfig() (c Config) {
	c.Prepare()
	return c
}

func (c *Config) Prepare(preferred, fallback string) {
	if preferred != nil {
		c.Name = preferred
		return
	}

	// Use the env variable if it exists
	if name, exists := os.LookupEnv(FullnameEnvVariable); exists {
		c.Name = name
		return
	}

	// Attempt to use local gitconfig for name
	if name, err := gitconfig.Username(); err != nil {
		c.Name = name
		return
	}

	// Attempt to use global gitconfig for name
	if name, err := gitconfig.Global("user.name"); err != nil {
		c.Name = name
		return
	}

	// Use the non-nil fallback
	if fallback != nil {
		c.Name = fallback
		return
	}

	// Finally use the default
	c.Name = defaultName
}
