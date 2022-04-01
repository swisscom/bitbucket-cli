package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SecurityCmd struct {
	Scan *SecurityScanCmd `arg:"subcommand:scan"`
}

type SecurityScanCmd struct {
}

func (b *BitbucketCLI) securityCmd(cmd *RepoCmd) {
	if cmd.SecurityCmd.Scan != nil {
		// Do scan
		b.triggerRepoScan(cmd.ProjectKey, cmd.Slug)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}

func (b *BitbucketCLI) triggerRepoScan(projectKey string, slug string) {
	b.logger.Debugf("triggering repo scan for %s/%s", projectKey, slug)

	// Custom endpoint!
	// https://docs.soteri.io/security-for-bitbucket/3.17.0/(3.17.0)-REST-API-for-Scripting-&-Automation.14602141697.html#id-(3.17.0)RESTAPIforScripting&Automation-Kickingoffanewrepositoryscan

	triggerScanUrl, err := b.restUrl.Parse(
		fmt.Sprintf("security/1.0/scan/%s/repos/%s",
			url.PathEscape(projectKey),
			url.PathEscape(slug),
		),
	)

	if err != nil {
		b.logger.Fatalf("unable to parse url: %v", err)
	}

	b.logger.Debugf("POST %s", triggerScanUrl.String())

	req, err := http.NewRequest(http.MethodPost, triggerScanUrl.String(), nil)
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

	b.logger.Infof("Scan successfully triggered")
}
