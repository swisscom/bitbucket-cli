package cli

type PrCmd struct {
	Create *PrCreateCmd `arg:"subcommand:create"`
	List   *PrListCmd   `arg:"subcommand:list"`
}

func (b *BitbucketCLI) RunPRCmd(cmd *PrCmd) {
	if cmd == nil {
		return
	}
	if cmd.Create != nil {
		b.RunPRCreateCmd(cmd.Create)
		return
	}
	if cmd.List != nil {
		b.RunPRListCmd(cmd.List)
		return
	}
	b.logger.Fatal(errSpecifySubcommand)
}
