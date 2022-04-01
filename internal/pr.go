package cli

type PrCmd struct {
	List *PrListCmd `arg:"subcommand:list"`
}

func (b *BitbucketCLI) RunPRCmd(cmd *PrCmd) {
	if cmd == nil {
		return
	}
	if cmd.List != nil {
		b.RunPRListCmd(cmd.List)
		return
	}
	b.logger.Fatal(errSpecifySubcommand)
}
