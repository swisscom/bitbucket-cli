package cli

import (
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"regexp"
	"strings"
)

type RepoBranchListCmd struct {
	Filter string `arg:"-f,--filter" help:"Filter to match branch names against (contains)"`
	Prefix string `arg:"-p,--prefix" help:"Only list branches that start with this prefix"`
	Regex  string `arg:"-r,--regex" help:"Only list branches that start with this prefix"`
}

func (b *BitbucketCLI) branchCmdList(cmd *RepoCmd) {
	if cmd == nil || cmd.BranchCmd == nil || cmd.BranchCmd.List == nil {
		return
	}

	filterFunction := func(branch bitbucketv1.Branch) bool { return true }

	list := cmd.BranchCmd.List
	if list.Prefix != "" {
		prevFilter := filterFunction
		filterFunction = func(branch bitbucketv1.Branch) bool {
			return prevFilter(branch) && strings.HasPrefix(branch.DisplayID, list.Prefix)
		}
	}

	if list.Regex != "" {
		regex, err := regexp.Compile(list.Regex)

		if err != nil {
			b.logger.Warnf("Regex %s is not valid, will be skipped", list.Regex)
		} else {
			prevFilter := filterFunction
			filterFunction = func(branch bitbucketv1.Branch) bool {
				return prevFilter(branch) && regex.MatchString(branch.DisplayID)
			}
		}
	}
	optionals := make(map[string]interface{})
	if list.Filter != "" {
		optionals["filterText"] = list.Filter
	}
	response, err := b.client.DefaultApi.GetBranches(cmd.ProjectKey, cmd.Slug, optionals)

	if err != nil {
		b.logger.Fatalf("Failed to fetch branches %s", err.Error())
		return
	}
	branches, err := bitbucketv1.GetBranchesResponse(response)

	if err != nil {
		b.logger.Fatalf("Failed to parse branches response %s", err.Error())
		return
	}

	for _, branch := range branches {
		if filterFunction(branch) {
			fmt.Printf("%s \n", branch.DisplayID)
		}
	}

}
