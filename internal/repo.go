package cli

type RepoCmd struct {
	ProjectKey string `arg:"-k,--key,required" help:"Project Key (e.g: TOOL)"`
	Slug       string `arg:"-n,--name,required" help:"Slug of the repository"`

	PrCmd *RepoPrCmd `arg:"subcommand:pr"`
}

func (b *BitbucketCLI) RunRepoCmd(cmd *RepoCmd) {
	if cmd == nil {
		return
	}

	if cmd.PrCmd != nil {
		b.repoPrCmd(cmd)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}
