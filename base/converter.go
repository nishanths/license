package base

import (
	"encoding/json"
	"regexp"
	"strings"
)

var placeholders = map[string]string{ // depends on api.github.com format
	"year":     "Year",
	"fullname": "Name",
}

var placeholdersRx *regexp.Regexp

func init() {
	var keys []string
	for key := range placeholders {
		keys = append(keys, key)
	}
	placeholdersRx = regexp.MustCompile("\\[(" + strings.Join(keys, "|") + ")\\]")
}

func textTemplateString(l *License) string {
	return placeholdersRx.ReplaceAllStringFunc(l.Body, func(m string) string {
		if s := placeholdersRx.FindStringSubmatch(m); s != nil && len(s) > 0 {
			k := s[1]
			if v, exists := placeholders[k]; exists {
				return "{{." + v + "}}"
			}
		}

		return m
	})
}

func jsonToList(content []byte) ([]License, error) {
	var licenses []License
	if err := json.Unmarshal(content, &licenses); err != nil {
		return nil, err
	}
	return licenses, nil
}

func jsonToLicense(content []byte) (License, error) {
	var full License
	if err := json.Unmarshal(content, &full); err != nil {
		return full, err
	}
	return full, nil
}
