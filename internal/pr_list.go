package cli

import (
	"encoding/json"
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"regexp"
	"strings"
)

type PrListCmd struct {
	State  string `arg:"-s,--state"`
	Output string `arg:"-o,--output"`

	FilterTitleRegex string `arg:"-t,--filter-title"`
	FilterDescRegex  string `arg:"-d,--filter-desc"`
}

/*
	/rest/api/1.0/dashboard/pull-requests?limit=100&state=OPEN
*/

func (b *BitbucketCLI) RunPRListCmd(cmd *PrListCmd) {
	if cmd == nil {
		return
	}

	var filterTitle *regexp.Regexp
	var filterDesc *regexp.Regexp
	var err error
	if cmd.FilterTitleRegex != "" {
		filterTitle, err = regexp.Compile(cmd.FilterTitleRegex)
		if err != nil {
			b.logger.Fatalf("unable to compile regex for title filtering: %v", err)
		}
	}

	if cmd.FilterDescRegex != "" {
		filterDesc, err = regexp.Compile(cmd.FilterDescRegex)
		if err != nil {
			b.logger.Fatalf("unable to compile regex for description filtering: %v", err)
		}
	}

	options := map[string]interface{}{}
	if cmd.State != "" {
		options["state"] = strings.ToUpper(cmd.State)
	}

	res, err := b.client.DefaultApi.GetPullRequests(options)
	if err != nil {
		b.logger.Fatalf("unable to list pull requests: %v", err)
		return
	}

	pr, err := bitbucket.GetPullRequestsResponse(res)
	if err != nil {
		b.logger.Fatalf("unable to parse PRs list: %v", err)
		return
	}

	var filteredPrs []bitbucket.PullRequest
	for _, v := range pr {
		if filterTitle != nil {
			if !filterTitle.MatchString(v.Title) {
				continue
			}
		}

		if filterDesc != nil {
			if !filterDesc.MatchString(v.Description) {
				continue
			}
		}

		filteredPrs = append(filteredPrs, v)
	}

	switch strings.ToLower(cmd.Output) {
	case "json":
		jsonOutput, err := json.Marshal(pr)
		if err != nil {
			b.logger.Fatalf("unable to marshal JSON: %v", err)
		}
		fmt.Printf("%s", jsonOutput)
	default:
		for _, v := range pr {
			fmt.Printf("%s - %s - %s\n", v.Title, v.Author.User.DisplayName, v.Links.Self[0].Href)
		}
	}
}
