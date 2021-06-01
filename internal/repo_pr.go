package cli

type RepoPrCmd struct {
	List   *RepoPrListCmd `arg:"subcommand:list"`
	Create *RepoPrCreateCmd `arg:"subcommand:create"`
}

func (b *BitbucketCLI) repoPrCmd(cmd *RepoCmd){
	if cmd == nil || cmd.PrCmd == nil {
		return
	}

	prCmd := cmd.PrCmd

	if prCmd.List != nil {
		b.repoPrList(cmd)
		return
	}

	if prCmd.Create != nil {
		b.repoPrCreate(cmd)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}