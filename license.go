package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

var licenses = map[string]struct {
	longName string
	template string
}{
	"afl-3.0":      {"Academic Free License", Afl30Template},
	"agpl-3.0":     {"GNU Affero General Public License v3.0", Agpl30Template},
	"al-2.0":       {"Artistic License 2.0", Al20Template},
	"apache-2.0":   {"Apache License 2.0", Apache20Template},
	"bsd-0-clause": {"BSD Zero Clause License", Bsd0ClauseTemplate},
	"bsd-2-clause": {"BSD 2-Clause \"Simplified\" License", Bsd2ClauseTemplate},
	"bsd-3-clause": {"BSD 3-Clause \"New\" or \"Revised\" License", Bsd3ClauseTemplate},
	"bsd-4-clause": {"BSD 4-clause \"Original\" or \"Old\" License", Bsd4ClauseTemplate},
	"cc0-1.0":      {"Creative Commons Zero v1.0 Universal", Cc010Template},
	"cc-by-4.0":    {"CC-BY 4.0 International Public License", Ccby4Template},
	"cc-by-sa-4.0": {"CC-BY-SA 4.0 International Public License", Ccbysa4Template},
	"ecl-2.0":      {"Educational Community License 2.0", Ecl20Template},
	"epl-2.0":      {"Eclipse Public License 2.0", Epl20Template},
	"eupl-1.2":     {"European Union Public Licence 1.2", Eupl12Template},
	"free-art-1.3": {"Free Art License 1.3", FreeArt13Template},
	"gpl-2.0":      {"GNU General Public License v2.0", Gpl20Template},
	"gpl-3.0":      {"GNU General Public License v3.0", Gpl30Template},
	"isc":          {"ISC License", ISCTemplate},
	"lgpl-2.1":     {"GNU Lesser General Public License v2.1", Lgpl21Template},
	"lgpl-3.0":     {"GNU Lesser General Public License v3.0", Lgpl30Template},
	"lppl":         {"LaTeX Project Public License", LpplTemplate},
	"odbl-1.0":     {"ODC Open Database License 1.0", Odbl10Template},
	"ofl-1.1":      {"SIL OPEN FONT LICENSE Version 1.1", Ofl11Template},
	"osl-3.0":      {"Open Software License v3.0", Osl30Template},
	"mit":          {"MIT License", MitTemplate},
	"mit-0":        {"MIT No Attribution", Mit0Template},
	"mpl-2.0":      {"Mozilla Public License 2.0", Mpl20Template},
	"ms-pl":        {"Microsoft Public License", MsplTemplate},
	"ms-rl":        {"Microsoft Reciprocal License", MsrlTemplate},
	"mulanpsl-2.0": {"木兰宽松许可证 第2版", Mulanpsl20Template},
	"unlicense":    {"The Unlicense", UnlicenseTemplate},
	"wtfpl":        {"Do What The Fuck You Want To Public License", WtfplTemplate},
	"zlib":         {"zlib License", ZlibTemplate},
}

func printList() {
	keys := make([]string, 0, len(licenses))

	for key := range licenses {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		stdout.Printf("%-14s(%s)", key, licenses[key].longName)
	}
}

func printLicense(license, output, name, year, project string) {
	file, ok := licenses[license]
	if !ok {
		stderr.Printf("unknown license %q\nrun \"license -list\" for list of available licenses", license)
		os.Exit(2)
	}

	t, err := template.New("license").Parse(file.template)
	if err != nil {
		stderr.Printf("internal: failed to parse license template for %s", license)
		os.Exit(1)
	}

	var outFile io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(filepath.Clean(output))
		if err != nil {
			stderr.Printf("failed to create file %s: %s", output, err)
			os.Exit(1)
		}
		outFile = f
	}

	if err := t.Execute(outFile, struct {
		Name    string
		Year    string
		Project string
	}{name, year, project}); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
