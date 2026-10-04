package cli

type RepoCmd struct {
	ProjectKey string `arg:"-k,--key,env:BITBUCKET_PROJECT" help:"Project key (e.g: TOOL); detected from git remote if not specified"`
	Slug       string `arg:"-n,--name,env:BITBUCKET_REPO" help:"Slug of the repository; detected from git remote if not specified (not used by \"create\")"`

	PrCmd       *RepoPrCmd     `arg:"subcommand:pr"`
	BranchCmd   *BranchCmd     `arg:"subcommand:branch"`
	SecurityCmd *SecurityCmd   `arg:"subcommand:security"`
	Get         *RepoGetCmd    `arg:"subcommand:get" help:"Show a repository"`
	Create      *RepoCreateCmd `arg:"subcommand:create" help:"Create a repository"`
	Update      *RepoUpdateCmd `arg:"subcommand:update" help:"Update a repository"`
	Delete      *RepoDeleteCmd `arg:"subcommand:delete" help:"Delete a repository"`
}

const errSlugRequired = "repository slug is required: use -n or run from inside a Bitbucket git repository."

// requireSlug stops the program unless a repository slug is known, either
// from --name or detected from the git remote.  The parser cannot enforce
// --name itself: detection may supply it, and "create" names the new
// repository with --display-name instead of using a slug at all.
func (b *BitbucketCLI) requireSlug(cmd *RepoCmd) {
	if cmd.Slug == "" {
		b.logger.Fatal(errSlugRequired)
	}
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

	// "create" needs only the project key; every other subcommand also needs
	// the slug of an existing repository.
	if cmd.Create != nil {
		b.repoCreate(cmd)
		return
	}

	if cmd.PrCmd != nil {
		b.requireSlug(cmd)
		b.repoPrCmd(cmd)
		return
	}

	if cmd.BranchCmd != nil {
		b.requireSlug(cmd)
		b.branchCmd(cmd)
		return
	}

	if cmd.SecurityCmd != nil {
		b.requireSlug(cmd)
		b.securityCmd(cmd)
		return
	}

	if cmd.Get != nil {
		b.requireSlug(cmd)
		b.repoGet(cmd)
		return
	}

	if cmd.Update != nil {
		b.requireSlug(cmd)
		b.repoUpdate(cmd)
		return
	}

	if cmd.Delete != nil {
		b.requireSlug(cmd)
		b.repoDelete(cmd)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}
