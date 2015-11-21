package remote

import (
	"encoding/json"
	// "fmt"
	"github.com/google/go-querystring/query"
	"github.com/nishanths/license/base"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	GitHubAPIBaseURL      = "https://api.github.com"
	GitHubAPILicensesPath = "/licenses"
	GitHubAPIAccept       = "application/vnd.github.drax-preview+json application/vnd.github.v3+json"
)

type option struct {
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
}

func do(req *http.Request) ([]byte, error) {
	client := http.Client{}

	req.Header.Add("Accept", GitHubAPIAccept)

	v, err := query.Values(option{os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET")})
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = v.Encode()

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

func Info(l *base.License) (base.License, error) {
	var full base.License

	req, err := http.NewRequest("GET", l.URL, nil)
	if err != nil {
		return full, err
	}

	body, err := do(req)
	if err := json.Unmarshal(body, &full); err != nil {
		return full, err
	}

	return full, nil
}
