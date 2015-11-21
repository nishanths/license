package base

import (
	"fmt"
)

type License struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	URL            string   `json:"url"`
	HtmlURL        string   `json:"html_url"`
	Featured       bool     `json:"featured"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Implementation string   `json:"implementation"`
	Required       []string `json:"required"`
	Permitted      []string `json:"permitted"`
	Forbidden      []string `json:"forbidden"`
	Body           string   `json:"body"`
}

func (l *License) String() string {
	return fmt.Sprintf("Key: %s, Name: %s", l.Key, l.Name)
}

func (l *License) TextTemplate() string {
	return PlaceholdersRx.ReplaceAllStringFunc(l.Body, func(m string) string {
		if s := PlaceholdersRx.FindStringSubmatch(m); s != nil && len(s) > 0 {
			k := s[1]
			return "{{." + Placeholders[k] + "}}"
		}

		return m
	})
}
