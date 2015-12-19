package base

import (
	"fmt"
)

type License struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	Url            string   `json:"url"`
	HtmlUrl        string   `json:"html_url"`
	Featured       bool     `json:"featured"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Implementation string   `json:"implementation"`
	Required       []string `json:"required"`
	Permitted      []string `json:"permitted"`
	Forbidden      []string `json:"forbidden"`
	Body           string   `json:"body"`
}

type ByLicenseKey []License

func (a ByLicenseKey) Len() int {
	return len(a)
}

func (a ByLicenseKey) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByLicenseKey) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

func (l *License) String() string {
	return fmt.Sprintf("{Key: %s, Name: %s}", l.Key, l.Name)
}
