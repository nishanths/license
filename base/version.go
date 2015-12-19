package base

import "fmt"

// Version prints the version number
func Version() error {
	fmt.Println(applicationVersion)
	return nil
}
