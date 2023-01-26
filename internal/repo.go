package cli

type RepoCmd struct {
	ProjectKey string `arg:"-k,--key,required,env:BITBUCKET_PROJECT" help:"Project AccessToken (e.g: TOOL)"`
	Slug       string `arg:"-n,--name,required,env:BITBUCKET_REPO" help:"Slug of the repository"`

	PrCmd       *RepoPrCmd   `arg:"subcommand:pr"`
	BranchCmd   *BranchCmd   `arg:"subcommand:branch"`
	SecurityCmd *SecurityCmd `arg:"subcommand:security"`
}

func (b *BitbucketCLI) RunRepoCmd(cmd *RepoCmd) {
	if cmd == nil {
		return
	}

	if cmd.PrCmd != nil {
		b.repoPrCmd(cmd)
		return
	}

	if cmd.BranchCmd != nil {
		b.branchCmd(cmd)
		return
	}

	if cmd.SecurityCmd != nil {
		b.securityCmd(cmd)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}
