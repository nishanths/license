// hgconfig provides methods to read Mercurial config items by name.
// It returns the results normally obtained from running "hg config [name]"
package hgconfig

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// future
const (
	configItemSeparator = "\n"
	nameValueSeparator  = "="
	nameSeparator       = "."
)

// ErrNameDoesNotExist is returned when a queried name is not found
// by "hg config"
type ErrNameDoesNotExist struct {
	Name string
}

func (err *ErrNameDoesNotExist) Error() string {
	return fmt.Sprintf("the name \"%s\" does not exist", err.Name)
}

func execHg(args []string) (string, error) {
	cmd := exec.Command("hg", args...)
	output, err := cmd.Output()

	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 { // `hg config` returns 1 if name does not exist
				return "", &ErrNameDoesNotExist{Name: args[len(args)-1]}
			}
		}
		return "", err
	}

	return strings.Trim(string(output), "\n"), nil
}

func execHgConfig(args ...string) (string, error) {
	return execHg(append([]string{"config"}, args...))
}

// Get lets you read a hg config value by name.
// It returns the value returned by "hg config [name]" for the given name.
// If there was an error finding the name, a non-nil error is returned.
func Get(name string) (string, error) {
	return execHgConfig(name)
}

// Username is a convenience function for getting "ui.username".
// This is the same as calling Get("ui.username").
func Username() (string, error) {
	return Get("ui.username")
}
