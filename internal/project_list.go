package cli

import (
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/sirupsen/logrus"
	"net/url"
)

type ProjectListCmd struct {
}

func (b *BitbucketCLI) projectList(cmd *ProjectCmd) {
	if cmd == nil {
		return
	}

	var repositories []bitbucket.Repository
	start := 0

	// Fetch all repos
	for {
		// List project repositories
		res, err := b.client.DefaultApi.GetRepositoriesWithOptions(cmd.Key,
			map[string]interface{}{
				"start": start,
			},
		)
		if err != nil {
			logrus.Fatal(err)
		}
		pageRepos, err := bitbucket.GetRepositoriesResponse(res)
		if err != nil {
			logrus.Fatalf("unable to parse repositories response: %v", err)
		}
		repositories = append(repositories, pageRepos...)
		hasNextPage, nextPageStart := bitbucket.HasNextPage(res)
		if !hasNextPage {
			break
		}
		start = nextPageStart
	}

	for _, v := range repositories {
		// Get HTTP Clone URL
		var cloneUrl = ""
		for _, cUrl := range v.Links.Clone {
			mUrl, err := url.Parse(cUrl.Href)
			if err != nil {
				continue
			}

			if mUrl.Scheme == "https" {
				cloneUrl = mUrl.String()
				break
			}
		}
		fmt.Printf("%s\t%s\n", v.Name, cloneUrl)
	}
}
