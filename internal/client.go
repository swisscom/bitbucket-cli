package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Custom Client to perform custom REST requests

func (b *BitbucketCLI) doReq(req *http.Request) (*http.Response, error) {
	err := b.prepareRequest(req)
	if err != nil {
		return nil, err
	}
	return b.httpClient.Do(req)
}

func (b *BitbucketCLI) getUrl() *url.URL {
	return b.restUrl
}

// apiBaseUrl returns the REST base URL without a trailing slash, e.g.
// https://git.example.com/rest, ready to have an API path appended to it.
func (b *BitbucketCLI) apiBaseUrl() string {
	return strings.TrimRight(b.restUrl.String(), "/")
}

func (b *BitbucketCLI) prepareRequest(req *http.Request) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	b.auth.AddHeaders(req)

	// https://confluence.atlassian.com/cloudkb/xsrf-check-failed-when-calling-cloud-apis-826874382.html
	req.Header.Set("X-Atlassian-Token", "no-check")
	return nil
}
