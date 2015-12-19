package base

import "fmt"

type helpLine struct {
	Left  string
	Right string
}

func (l *helpLine) String() string {
	return fmt.Sprintf("%s%-14s%s", indent, l.Left, l.Right)
}

func printCommands() {
	fmt.Println("Additional commands:")
	for _, c := range []helpLine{
		{"ls", "list locally available license names"},
		{"ls-remote", "list remote license names"},
		{"update", "update local licenses to latest remote versions"},
		{"help", "show help information"},
		{"version", "print current version"},
	} {
		fmt.Println(&c)
	}
}

func printExamples() {
	fmt.Println("Examples:")
	for _, c := range []helpLine{
		{"license mit", ""},
		{"license -o LICENSE.txt mit", ""},
		{"license -y 2013 -n Alice isc", ""},
	} {
		fmt.Println(&c)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(indent + "license [-y <year>] [-n <name>] [-o <filename>] <license-name>")
}

func printOptions() {
	fmt.Println("Options:")
	for _, c := range []helpLine{
		{"-y, --year", "year on the license"},
		{"-n, --name", "name on the license"},
		{"-o, --output", "filename to save license"},
	} {
		fmt.Println(&c)
	}
}

// Help prints help information
// for the program to the console
func Help() error {
	// Heading
	fmt.Println("Command-line license generator.")
	fmt.Println()

	// Usage
	printUsage()
	fmt.Println()

	// Example
	printExamples()
	fmt.Println()

	// Additional commands
	printCommands()
	fmt.Println()

	// Note
	fmt.Println("Run \"license ls\" to see list of available license names.")

	return nil
}
