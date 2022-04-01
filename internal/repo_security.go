package cli

import (
	"fmt"
	"github.com/fatih/color"
)

type SecurityCmd struct {
	Scan   *SecurityScanCmd   `arg:"subcommand:scan"`
	Result *SecurityResultCmd `arg:"subcommand:result"`
}

func printGreenBackground(text string) {
	c := color.New(color.FgBlack).Add(color.BgGreen)
	_, _ = c.Printf(" %s ", text)
}

func printRedBackground(text string) {
	c := color.New(color.FgBlack).Add(color.BgRed)
	_, _ = c.Printf(" %s ", text)
}

func (b *BitbucketCLI) securityCmd(cmd *RepoCmd) {
	if cmd.SecurityCmd.Scan != nil {
		// Do scan
		b.triggerRepoScan(cmd.ProjectKey, cmd.Slug)
		return
	}

	if cmd.SecurityCmd.Result != nil {
		scanResult := b.getScanResult(cmd.ProjectKey, cmd.Slug)
		if scanResult.Total == 0 {
			printGreenBackground(fmt.Sprintf(
				"No issues found in %s/%s", cmd.ProjectKey, cmd.Slug),
			)
		} else {
			printRedBackground(
				fmt.Sprintf("%d issues found in %s/%s", scanResult.Total, cmd.ProjectKey, cmd.Slug),
			)
		}
		fmt.Printf("\n")
		return
	}

	b.logger.Fatal(errSpecifySubcommand)
}
