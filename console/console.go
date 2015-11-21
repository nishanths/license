package console

import (
	"fmt"
	"github.com/fatih/color"
)

func Error(s string) {
	fmt.Println(fmt.Sprintf("%s %s", color.RedString("ERR "), s))
}

func Info(s string) {
	fmt.Println(fmt.Sprintf("%s %s", color.YellowString("INFO"), s))
}
