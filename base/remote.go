package base

import (
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	gitHubClientIdEnvVariable     = "GITHUB_CLIENT_ID"
	gitHubClientSecretEnvVariable = "GITHUB_CLIENT_SECRET"
	gitHubAPIBaseURL              = "https://api.github.com"
	gitHubAPILicensesPath         = "/licenses"
	gitHubAPIAccept               = "application/vnd.github.drax-preview+json application/vnd.github.v3+json"
)

// fetch performs a HTTP request after appending required headers,
// and returns the response bytes and an error, if any.
func fetch(req *http.Request) ([]byte, error) {
	type option struct {
		ClientID     string `url:"client_id"`
		ClientSecret string `url:"client_secret"`
	}

	client := &http.Client{}

	// additional http headers
	req.Header.Add("Accept", gitHubAPIAccept)

	// query string
	queryOpt := option{os.Getenv(gitHubClientIdEnvVariable), os.Getenv(gitHubClientSecretEnvVariable)}
	queryValues, err := query.Values(queryOpt)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = queryValues.Encode()

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func fetchIndex() ([]byte, error) {
	req, err := http.NewRequest("GET", gitHubAPIBaseURL+gitHubAPILicensesPath, nil)

	if err != nil {
		return nil, err
	}

	return fetch(req)
}

func (l *License) fetchFullInfo() ([]byte, error) {
	req, err := http.NewRequest("GET", l.Url, nil)

	if err != nil {
		return nil, err
	}

	return fetch(req)
}
