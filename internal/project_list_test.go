package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"github.com/swisscom/bitbucket-cli/test"
	"testing"
)

func TestProjectList(t *testing.T) {
	c := test.MustGetCLI()
	c.RunProjectCmd(&cli.ProjectCmd{
		Key:   "TOOL",
		List:  &cli.ProjectListCmd{},
		Clone: nil,
	})
}
