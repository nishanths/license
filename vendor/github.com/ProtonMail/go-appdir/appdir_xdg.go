// +build !darwin,!windows

package appdir

import (
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) UserConfig() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".config")
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		baseDir = d
	}

	return filepath.Join(baseDir, d.name)
}

func (d *dirs) UserCache() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".cache")
	if d := os.Getenv("XDG_CACHE_HOME"); d != "" {
		baseDir = d
	}

	return filepath.Join(baseDir, d.name)
}

func (d *dirs) UserLogs() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".local", "state")
	if d := os.Getenv("XDG_STATE_HOME"); d != "" {
		baseDir = d
	}

	return filepath.Join(baseDir, d.name)
}

func (d *dirs) UserData() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".local", "share")
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		baseDir = d
	}

	return filepath.Join(baseDir, d.name)
}
