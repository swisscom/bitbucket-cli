package cli

import (
	"fmt"
)

type RepoPrApproveCmd struct {
	Id int64 `arg:"-i,--id" help:"id of the PR"`
}

func (b *BitbucketCLI) repoPrApprove(cmd *RepoCmd) {
	if cmd == nil || cmd.PrCmd == nil || cmd.PrCmd.Approve == nil {
		return
	}
	approve := cmd.PrCmd.Approve

	if _, err := b.client.DefaultApi.Approve(cmd.ProjectKey, cmd.Slug, approve.Id); err != nil {
		b.logger.Fatalf("unable to approve PR: %v", err)
	} else {
		fmt.Printf("PR %d succesfully approved", approve.Id)
	}
}
