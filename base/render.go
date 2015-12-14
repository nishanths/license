package base

import (
	"fmt"
	"io"
	"text/template"
	"time"
)

const indent = "  "

type Option struct {
	Year int
	Name string
}

func NewOption(n string) (o Option) {
	o.Year = time.Now().Year()
	o.Name = n
	return
}

func singleFormatString(l *License) string {
	return fmt.Sprintf("%s%-15s(%s)", indent+indent, l.Key, l.Name)
}

func RenderList(licenses *[]License) {
	for _, l := range *licenses {
		fmt.Println(singleFormatString(&l))
	}
}

func RenderTemplate(t *template.Template, o *Option, w io.Writer) (err error) {
	err = t.ExecuteTemplate(w, t.Name(), o)
	return
}
