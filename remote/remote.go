package remote

import (
	"encoding/json"
	// "fmt"
	"github.com/nishanths/license/base"
	"io/ioutil"
	"net/http"
)

const (
	GH_API_BASE        = "https://api.github.com"
	GH_LICENSES_PATH   = "/licenses"
	GH_LICENSES_ACCEPT = "application/vnd.github.drax-preview+json application/vnd.github.v3+json"
)

type LicensesResult struct {
	Licenses *[]base.License
	Error    error
}

func FetchLicensesList(ch chan *LicensesResult) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", GH_API_BASE+GH_LICENSES_PATH, nil)
	if err != nil {
		ch <- &LicensesResult{nil, err}
	}

	req.Header.Add("Accept", GH_LICENSES_ACCEPT)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		ch <- &LicensesResult{nil, err}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- &LicensesResult{nil, err}
	}

	licenses := make([]base.License, 0)
	json.Unmarshal(body, &licenses)
	ch <- &LicensesResult{&licenses, err}
}

func UpdateLicense() {
	// fetchCh := make(chan *LicensesResult)

	// go FetchLicensesList(fetchCh)

	// lr := <-fetchCh
	// if lr.Error != nil {
	// 	return
	// }

	// fmt.Println(*lr.Licenses)
}
