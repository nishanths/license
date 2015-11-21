package base

import (
	"regexp"
	"strings"
)

var Placeholders = map[string]string{
	"year":     "Year",
	"fullname": "Name",
}

var PlaceholdersRx *regexp.Regexp

func init() {
	keys := make([]string, 0)
	for key := range Placeholders {
		keys = append(keys, key)
	}
	PlaceholdersRx = regexp.MustCompile("\\[(" + strings.Join(keys, "|") + ")\\]")
}
