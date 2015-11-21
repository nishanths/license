package base

import (
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"
)

type Option struct {
	Year string
	Name string
}

func NewOption(n string) (o Option) {
	o.Year = strconv.Itoa(time.Now().Year())
	o.Name = n
	return
}

func singleFormatString(l *License) string {
	return fmt.Sprintf("* %s (%s)", l.Key, l.Name)
}

func RenderList(licenses *[]License) {
	for _, l := range *licenses {
		fmt.Println(singleFormatString(&l))
	}
}

func RenderTemplate(t *template.Template, o *Option) (err error) {
	err = t.ExecuteTemplate(os.Stdout, t.Name(), o)
	return
}
