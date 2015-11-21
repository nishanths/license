package console

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Error(s string) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s %s", color.RedString("ERR "), s))
}

func Info(s string) {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("%s %s", color.YellowString("INFO"), s))
}
