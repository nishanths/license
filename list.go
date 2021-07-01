package main

var (
	licensesList = map[string]struct {
		longName string
		template string
	}{
		"agpl-3.0":     {"GNU Affero General Public License v3.0", Agpl30Template},
		"apache-2.0":   {"Apache License 2.0", Apache20Template},
		"bsd-2-clause": {"BSD 2-Clause \"Simplified\" License", Bsd2ClauseTemplate},
		"bsd-3-clause": {"BSD 3-Clause \"New\" or \"Revised\" License", Bsd3ClauseTemplate},
		"cc0-1.0":      {"Creative Commons Zero v1.0 Universal", Cc010Template},
		"epl-2.0":      {"Eclipse Public License 2.0", Epl20Template},
		"free-art-1.3": {"Free Art License 1.3", FreeArt13Template},
		"gpl-2.0":      {"GNU General Public License v2.0", Gpl20Template},
		"gpl-3.0":      {"GNU General Public License v3.0", Gpl30Template},
		"lgpl-2.1":     {"GNU Lesser General Public License v2.1", Lgpl21Template},
		"lgpl-3.0":     {"GNU Lesser General Public License v3.0", Lgpl30Template},
		"mit":          {"MIT License", MitTemplate},
		"mpl-2.0":      {"Mozilla Public License 2.0", Mpl20Template},
		"unlicense":    {"The Unlicense", UnlicenseTemplate},
		"wtfpl":        {"Do What The Fuck You Want To Public License", WtfplTemplate},
	}
)

func printList() {
	for key, license := range licensesList {
		stdout.Printf("%-14s(%s)", key, license.longName)
	}
}
