package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"github.com/swisscom/bitbucket-cli/test"
	"testing"
)

// TestRepoCRUD creates, reads, updates and deletes a throw-away repository
// against a real Bitbucket server.
func TestRepoCRUD(t *testing.T) {
	test.SkipIfNoServer(t)
	c := test.MustGetCLI()

	// No spaces or capitals, so the slug Bitbucket derives equals the name.
	name := "bitbucket-cli-crud-test"

	c.RunRepoCmd(&cli.RepoCmd{
		ProjectKey: "TOOL",
		Create: &cli.RepoCreateCmd{
			DisplayName: name,
			Description: "temporary repository created by the bitbucket-cli tests",
		},
	})
	c.RunRepoCmd(&cli.RepoCmd{
		ProjectKey: "TOOL",
		Slug:       name,
		Get:        &cli.RepoGetCmd{},
	})
	forkable := false
	c.RunRepoCmd(&cli.RepoCmd{
		ProjectKey: "TOOL",
		Slug:       name,
		Update:     &cli.RepoUpdateCmd{Forkable: &forkable},
	})
	c.RunRepoCmd(&cli.RepoCmd{
		ProjectKey: "TOOL",
		Slug:       name,
		Delete:     &cli.RepoDeleteCmd{Yes: true},
	})
}
