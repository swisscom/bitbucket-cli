package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SecurityResultCmd struct{}

type ScanResult struct {
	ScanKey      string        `json:"scanKey"`
	Scanned      bool          `json:"scanned"`
	Actual       bool          `json:"actual"`
	Progress     int           `json:"progress"`
	Running      bool          `json:"running"`
	Scheduled    bool          `json:"scheduled"`
	InvalidLines []interface{} `json:"invalidLines"`
	Total        int           `json:"total"`
}

func (b *BitbucketCLI) getScanResult(projectKey string, slug string) ScanResult {
	b.logger.Debugf("fetching security scan result for %s/%s", projectKey, slug)

	// Custom endpoint!
	// https://docs.soteri.io/security-for-bitbucket/3.17.0/(3.17.0)-REST-API-for-Scripting-&-Automation.14602141697.html#id-(3.17.0)RESTAPIforScripting&Automation-Fetchingscanresultsforaspecificbranch

	triggerScanUrl, err := b.restUrl.Parse(
		fmt.Sprintf("security/1.0/scan/%s/repos/%s",
			url.PathEscape(projectKey),
			url.PathEscape(slug),
		),
	)
	if err != nil {
		b.logger.Fatalf("unable to parse url: %v", err)
	}

	b.logger.Debugf("GET %s", triggerScanUrl.String())

	req, err := http.NewRequest(http.MethodGet, triggerScanUrl.String(), nil)
	if err != nil {
		b.logger.Fatalf("unable to create request: %v", err)
	}

	err = b.prepareRequest(req)
	if err != nil {
		b.logger.Fatalf("unable to prepare request: %v", err)
	}

	res, err := b.doReq(req)
	if err != nil {
		b.logger.Fatalf("unable to do request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		resBody, err := ioutil.ReadAll(res.Body)
		if err == nil {
			b.logger.Debugf("resp=%v", string(resBody))
		}
		b.logger.Fatalf("invalid status code received: 200 OK expected but %s received", res.Status)
	}

	// Decode response
	var scanResult ScanResult
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&scanResult)

	if err != nil {
		b.logger.Fatalf("unable to decode JSON: %v", err)
	}

	return scanResult
}
