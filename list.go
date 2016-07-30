package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/mitchellh/go-homedir"
	"github.com/nishanths/license/pkg/license"
)

// ByKey impelements sort.Interface to sort
// License by the Key field.
type ByKey []license.License

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// printList prints the supplied list of licenses to stdout.
// The list is sorted by the Key field before printing.
func printList(l []license.License) {
	sort.Sort(ByKey(l))
	logger.Print("Available licenses:\n")
	for _, l := range l {
		logger.Printf("  %-14s(%s)\n", l.Key, l.Name)
	}
}

// list prints a list of locally available licenses
// and exits.
func list() {
	h, err := homedir.Dir()
	if err != nil {
		errLogger.Println(err)
		os.Exit(1)
	}

	p := filepath.Join(h, ".license", "data", "licenses.json")
	if _, err := os.Stat(p); err != nil {
		errLogger.Println(err)
		os.Exit(1)
	}

	r, err := os.Open(p)
	if err != nil {
		errLogger.Println("failed to open licenses.json", err)
		os.Exit(1)
	}

	var lics []license.License
	if err := json.NewDecoder(r).Decode(&lics); err != nil {
		errLogger.Println("failed to decode licenses.json", err)
		os.Exit(1)
	}

	printList(lics)
	os.Exit(0)
}
