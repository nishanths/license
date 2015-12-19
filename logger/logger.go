package logger

import "fmt"

type logLevel struct {
	Verbose, Quiet bool
}

var globalLogLevel *logLevel

func init() {
	globalLogLevel = &logLevel{
		Verbose: false, Quiet: false,
	}
}

func (l *logLevel) outputAllowed() bool {
	return !l.Quiet
}

func (l *logLevel) verboseOutputAllowed() bool {
	return !l.Quiet && l.Verbose
}

func SetVerbose(b bool) {
	globalLogLevel.Verbose = b
}

func SetQuiet(b bool) {
	globalLogLevel.Quiet = b
}

// Print calls fmt.Print if quiet mode is off
func Print(args ...interface{}) {
	if globalLogLevel.outputAllowed() {
		fmt.Print(args...)
	}
}

// Printf calls fmt.Printf if quiet mode is off
func Printf(format string, args ...interface{}) {
	if globalLogLevel.outputAllowed() {
		fmt.Printf(format, args...)
	}
}

// Println calls fmt.Println if quiet mode is off
func Println(args ...interface{}) {
	if globalLogLevel.outputAllowed() {
		fmt.Println(args...)
	}
}

// VerbosePrint calls fmt.Print only when verbose logging is on
// and quiet mode is off
func VerbosePrint(args ...interface{}) {
	if globalLogLevel.verboseOutputAllowed() {
		fmt.Print(args...)
	}
}

// VerbosePrintf calls fmt.Printf only when verbose logging is on
// and quiet mode is off
func VerbosePrintf(format string, args ...interface{}) {
	if globalLogLevel.verboseOutputAllowed() {
		fmt.Printf(format, args...)
	}
}

// VerbosePrintln calls fmt.Println only when verbose logging is on
// and quiet mode is off
func VerbosePrintln(args ...interface{}) {
	if globalLogLevel.verboseOutputAllowed() {
		fmt.Println(args...)
	}
}
