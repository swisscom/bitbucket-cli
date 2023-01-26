package cli

type RepoPrCmd struct {
	Approve *RepoPrApproveCmd `arg:"subcommand:approve"`
	Create  *RepoPrCreateCmd  `arg:"subcommand:create"`
	List    *RepoPrListCmd    `arg:"subcommand:list"`
	Merge   *RepoPrMergeCmd   `arg:"subcommand:merge"`
}

func (b *BitbucketCLI) repoPrCmd(cmd *RepoCmd) {
	if cmd == nil || cmd.PrCmd == nil {
		return
	}

	prCmd := cmd.PrCmd

	if prCmd.Approve != nil {
		b.repoPrApprove(cmd)
		return
	} else if prCmd.Create != nil {
		b.repoPrCreate(cmd)
		return
	} else if prCmd.List != nil {
		b.repoPrList(cmd)
		return
	} else if prCmd.Merge != nil {
		b.repoPrMerge(cmd)
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}
