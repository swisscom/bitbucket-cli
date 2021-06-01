package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"github.com/swisscom/bitbucket-cli/test"
	"testing"
)

func TestProjectClone(t *testing.T) {
	c := test.GetCLI()
	c.RunProjectCmd(&cli.ProjectCmd{
		Key:  "TOOL",
		List: nil,
		Clone: &cli.ProjectCloneCmd{
			Output: "/tmp/git-repo",
			Branch: "master",
		},
	})
}
