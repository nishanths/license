package main

var (
	licensesList = []struct {
		key      string
		longName string
	}{
		{"agpl-3.0", "GNU Affero General Public License v3.0"},
		{"apache-2.0", "Apache License 2.0"},
		{"bsd-2-clause", "BSD 2-Clause \"Simplified\" License"},
		{"bsd-3-clause", "BSD 3-Clause \"New\" or \"Revised\" License"},
		{"cc0-1.0", "Creative Commons Zero v1.0 Universal"},
		{"epl-2.0", "Eclipse Public License 2.0"},
		{"free-art-1.3", "Free Art License 1.3"},
		{"gpl-2.0", "GNU General Public License v2.0"},
		{"gpl-3.0", "GNU General Public License v3.0"},
		{"lgpl-2.1", "GNU Lesser General Public License v2.1"},
		{"lgpl-3.0", "GNU Lesser General Public License v3.0"},
		{"mit", "MIT License"},
		{"mpl-2.0", "Mozilla Public License 2.0"},
		{"unlicense", "The Unlicense"},
		{"wtfpl", "Do What The Fuck You Want To Public License"},
	}
)

func printList() {
	for _, l := range licensesList {
		stdout.Printf("%-14s(%s)", l.key, l.longName)
	}
}
