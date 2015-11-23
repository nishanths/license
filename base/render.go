package base

import (
	"fmt"
	"math"
	"os"
	"text/template"
	"time"
)

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
	return fmt.Sprintf("   %s (%s)", l.Key, l.Name)
}

func RenderList(licenses *[]License, ms time.Duration) {
	for _, l := range *licenses {
		time.Sleep(ms * time.Duration(math.Pow(float64(10), float64(6))))
		fmt.Println(singleFormatString(&l))
	}
}

func RenderTemplate(t *template.Template, o *Option) (err error) {
	err = t.ExecuteTemplate(os.Stdout, t.Name(), o)
	return
}
