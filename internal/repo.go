package cli

type RepoCmd struct {
	ProjectKey string `arg:"-k,--key,env:BITBUCKET_PROJECT" help:"Project key (e.g: TOOL); detected from git remote if not specified"`
	Slug       string `arg:"-n,--name,env:BITBUCKET_REPO" help:"Slug of the repository; detected from git remote if not specified"`

	PrCmd       *RepoPrCmd   `arg:"subcommand:pr"`
	BranchCmd   *BranchCmd   `arg:"subcommand:branch"`
	SecurityCmd *SecurityCmd `arg:"subcommand:security"`
}

func (b *BitbucketCLI) RunRepoCmd(cmd *RepoCmd) {
	if cmd == nil {
		return
	}

	if cmd.ProjectKey == "" || cmd.Slug == "" {
		if ctx, err := GetRepoContext("."); err == nil {
			if cmd.ProjectKey == "" {
				cmd.ProjectKey = ctx.ProjectKey
			}
			if cmd.Slug == "" {
				cmd.Slug = ctx.Slug
			}
		}
	}

	if cmd.ProjectKey == "" {
		b.logger.Fatal("project key is required: use -k or run from inside a Bitbucket git repository.")
	}
	if cmd.Slug == "" {
		b.logger.Fatal("repository slug is required: use -n or run from inside a Bitbucket git repository.")
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
