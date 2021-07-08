package cli

import (
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"strings"
)

type RepoBranchCompareCmd struct {
	FromBranch string `arg:"-f,--from,required" help:"Name of the branch to be used as base"`
	ToBranch   string `arg:"-t,--to" help:"Name of the branch to be compared to the base"`
	Count      bool   `arg:"-c,--count" help:"Just output the number of commits between to and from branch"`
}

func (b *BitbucketCLI) branchCmdCompare(cmd *RepoCmd) {
	if cmd == nil || cmd.BranchCmd == nil || cmd.BranchCmd.Compare == nil {
		return
	}
	compare := cmd.BranchCmd.Compare

	response, err := b.client.DefaultApi.GetBranches(cmd.ProjectKey, cmd.Slug, nil)

	if err != nil {
		b.logger.Fatalf("Failed to fetch branches %s", err.Error())
		return
	}
	branches, err := bitbucketv1.GetBranchesResponse(response)

	fromBranch := findBranch(branches, cmd.BranchCmd.Compare.FromBranch)

	if fromBranch == nil {
		b.logger.Fatalf("Failed to find branch with name %s", compare.FromBranch)
		return
	}

	toBranch := findBranch(branches, cmd.BranchCmd.Compare.ToBranch)
	if toBranch == nil {
		b.logger.Fatalf("Failed to find branch with name %s", compare.FromBranch)
		return
	}

	optionals := make(map[string]interface{})
	if compare.Count {
		optionals["withCounts"] = true
	}

	optionals["since"] = fromBranch.LatestCommit
	optionals["until"] = toBranch.LatestCommit

	commitResponse, err := b.client.DefaultApi.GetCommits(cmd.ProjectKey, cmd.Slug, optionals)

	if err != nil {
		b.logger.Fatalf("Failed to fetch commits %s", err.Error())
		return
	}

	commits, err := bitbucketv1.GetCommitsResponse(commitResponse)

	if err != nil {
		b.logger.Fatalf("Failed to parse commits %s", err.Error())
		return
	}

	if compare.Count {
		fmt.Printf("%d", len(commits))
	} else {
		for _, c := range commits {
			fmt.Printf("%s %s \n", c.DisplayID, strings.Split(c.Message, "\n")[0])
		}
	}

}

func findBranch(branches []bitbucketv1.Branch, branchName string) *bitbucketv1.Branch {
	for _, b := range branches {
		if b.DisplayID == branchName {
			return &b
		}
	}

	return nil
}
