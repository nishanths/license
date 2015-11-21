package remote

import (
	"encoding/json"
	// "fmt"
	"github.com/nishanths/license/base"
	"io/ioutil"
	"net/http"
)

const (
	GitHubAPIBaseURL      = "https://api.github.com"
	GitHubAPILicensesPath = "/licenses"
	GitHubAPIAccept       = "application/vnd.github.drax-preview+json application/vnd.github.v3+json"
)

func do(req *http.Request) ([]byte, error) {
	client := http.Client{}
	req.Header.Add("Accept", GitHubAPIAccept)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func List() ([]base.License, error) {
	req, err := http.NewRequest("GET", GitHubAPIBaseURL+GitHubAPILicensesPath, nil)
	if err != nil {
		return nil, err
	}

	body, err := do(req)

	licenses := make([]base.License, 0)
	if err := json.Unmarshal(body, &licenses); err != nil {
		return nil, err
	}
	return licenses, nil
}

func Info(l *base.License) (*base.License, error) {
	req, err := http.NewRequest("GET", l.Url, nil)
	if err != nil {
		return nil, err
	}

	body, err := do(req)
	var full base.License
	if err := json.Unmarshal(body, &full); err != nil {
		return nil, err
	}

	return &full, nil
}
