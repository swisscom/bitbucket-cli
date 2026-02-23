package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"github.com/swisscom/bitbucket-cli/test"
	"testing"
)

func TestPRCreate(t *testing.T) {
	test.SkipIfNoServer(t)
	c := test.MustGetCLI()
	c.RunPRCmd(&cli.PrCmd{
		Create: &cli.PrCreateCmd{
			ProjectKey:  "TOOL",
			Slug:        "bitbucket-playground",
			Title:       "Test PR",
			Description: "Test PR created by bitbucket-cli",
			FromRef:     "refs/heads/feature/test",
			ToRef:       "refs/heads/master",
		},
	})
}
