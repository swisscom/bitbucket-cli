package cli

type BranchCmd struct {
	Compare *RepoBranchCompareCmd `arg:"subcommand:compare"`
	List    *RepoBranchListCmd    `arg:"subcommand:list"`
}

func (b *BitbucketCLI) branchCmd(cmd *RepoCmd) {
	if cmd.BranchCmd.Compare != nil {
		b.branchCmdCompare(cmd)
		return
	}
	if cmd.BranchCmd.List != nil {
		b.branchCmdList(cmd)
		return
	}
}
