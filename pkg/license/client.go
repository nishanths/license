// Package license implements a HTTP client for accessing the GitHub
// Licenses API. See:
//   https://developer.github.com/v3/licenses/
package license

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// BaseURL is the base URL of the Github Licenses API.
	BaseURL = "https://api.github.com/licenses"
)

var (
	// Header is headers included in each request.
	Header = map[string][]string{
		"Accept": []string{
			"application/vnd.github.drax-preview+json",
			"application/vnd.github.v3+json",
		},
	}
)

// StatusError is returned when the HTTP roundtrip succeeds, but the
// the response status does not equal http.StatusOK.
type StatusError struct {
	Status     string
	StatusCode int
	Details    struct {
		Message string `json:"message"`
		Errors  []struct {
			Resource string `json:"resource"`
			Field    string `json:"field"`
			Code     string `json:"code"`
		} `json:"errors"`
	}
}

func (e StatusError) Error() string {
	s := fmt.Sprintf("license: HTTP response error: %s", e.Status)
	if e.Details.Message != "" {
		s += ": " + e.Details.Message
	}
	return s
}

// Client represents a HTTP client for making requests to the
// GitHub Licenses API. Use NewClient to create a Client
// ready for use.
//
// If Token and Username are both set in Config, they are used for
// authentication. Otherwise, if ClientID and ClientSecret
// are both set, they are used for authentication. If none of the
// above fields are set, no authentication is used.
type Client struct {
	HTTPClient *http.Client
	Config
}

// NewClient returns a ready-to-use client to make calls to the
// GitHub Licenses API.
func NewClient() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		Config: Config{
			BaseURL: BaseURL,
		},
	}
}

// Config is configuration for Client.
type Config struct {
	// Token is GitHub personal API token for authentication.
	Token string

	// Username is GitHub username for authentication.
	Username string

	// ClientID is the GitHub client ID for the API.
	// See https://developer.github.com/v3/oauth/.
	// If either ClientID or ClientSecret is empty,
	// then ClientID and ClientSecret are not included in the request.
	ClientID string

	// ClientSecret is the GitHub client secret for the API.
	// See https://developer.github.com/v3/oauth/.
	ClientSecret string

	// BaseURL is the base URL for the API.
	// Useful in tests.
	BaseURL string

	// Header is a map of custom headers to add to outgoing request.
	// If nil or empty, no custom headers will be added.
	Header map[string][]string
}

func (c *Client) makeQuery(vals url.Values) string {
	// Shallow copy, but OK.
	query := url.Values{}
	for k, v := range vals {
		query[k] = v
	}

	if (c.Token == "" && c.Username == "") &&
		(c.ClientID != "" && c.ClientSecret != "") {
		query["client_id"] = []string{c.ClientID}
		query["client_secret"] = []string{c.ClientSecret}
	}

	return query.Encode()
}

func (c *Client) addHeaders(req *http.Request) {
	for k, v := range Header {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}
	for k, v := range c.Header {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}
}

func (c *Client) doReq(method, path string, body io.Reader, vals url.Values) (io.ReadCloser, error) {
	u := c.BaseURL + path
	q := c.makeQuery(vals)
	if q != "" {
		u += "?" + q
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	if c.Username != "" && c.Token != "" {
		req.SetBasicAuth(c.Username, c.Token)
	}

	c.addHeaders(req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		se := StatusError{
			Status:     res.Status,
			StatusCode: res.StatusCode,
		}
		json.NewDecoder(res.Body).Decode(&se.Details) // Ignore error.
		return nil, se
	}

	return res.Body, nil
}

// List returns the list of licenses.
// Only the following fields will be set in returned licenses.
//   Key, Name, URL, Featured
func (c *Client) List() ([]License, error) {
	data, err := c.doReq("GET", "", nil, nil)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var lics []License
	if err := json.NewDecoder(data).Decode(&lics); err != nil {
		return nil, err
	}
	return lics, err
}

// Info returns the license for key.
// Example keys are "mit", "lgpl-3.0", etc.
func (c *Client) Info(key string) (License, error) {
	data, err := c.InfoJSON(key)
	if err != nil {
		return License{}, err
	}
	defer data.Close()

	var lic License
	if err := json.NewDecoder(data).Decode(&lic); err != nil {
		return License{}, err
	}
	return lic, err
}

// InfoJSON returns the license as raw JSON for key.
// Example keys are "mit", "lgpl-3.0", etc.
// If the error is non-nil, the caller is reponsible
// for closing the returned io.ReadCloser.
func (c *Client) InfoJSON(key string) (io.ReadCloser, error) {
	return c.doReq("GET", "/"+key, nil, nil)
}
