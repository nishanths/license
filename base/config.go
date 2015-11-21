package base

import (
	gitconfig "github.com/tcnksm/go-gitconfig"
	"os"
)

const (
	FullnameEnvVariable = "LICENSE_FULLNAME"
	defaultFullname     = ""
)

type Config struct {
	Name string `json:"name"`
}

func (c *Config) Prepare(preferred, fallback string) {
	if preferred != "" {
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

	// Use the non-empty fallback
	if fallback != "" {
		c.Name = fallback
		return
	}

	// Finally use the default
	c.Name = defaultFullname
}
