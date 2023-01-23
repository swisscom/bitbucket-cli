package cli

import (
	"fmt"
	goBitBucket "github.com/gfleury/go-bitbucket-v1"
)

type RepoPrMergeCmd struct {
	Id int `arg:"-i,--id" help:"id of the PR"`
}

func (b *BitbucketCLI) repoPrMerge(cmd *RepoCmd) {
	if cmd == nil || cmd.PrCmd == nil || cmd.PrCmd.Merge == nil {
		return
	}
	merge := cmd.PrCmd.Merge
	if prResponse, err := b.client.DefaultApi.GetPullRequest(cmd.ProjectKey, cmd.Slug, merge.Id); err != nil {
		b.logger.Fatalf("unable to find PR: %v", err)
	} else if pr, err := goBitBucket.GetPullRequestResponse(prResponse); err != nil {
		b.logger.Fatalf("could not retrieve PR from response: %v", err)
	} else if !pr.Open {
		b.logger.Fatal("PR is in a closed state")
	} else {
		if mergeResponse, err := b.client.DefaultApi.Merge(
			cmd.ProjectKey,
			cmd.Slug,
			merge.Id,
			map[string]interface{}{"version": pr.Version},
			nil,
			nil); err != nil {
			b.logger.Fatalf("unable to approve PR: %v", err)
		} else {
			fmt.Println(mergeResponse.Message)
		}
	}
}
