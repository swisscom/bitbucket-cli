package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"github.com/swisscom/bitbucket-cli/test"
	"testing"
)

func TestPRList(t *testing.T) {
	test.SkipIfNoServer(t)
	c := test.MustGetCLI()
	c.RunPRCmd(&cli.PrCmd{
		List: &cli.PrListCmd{},
	})
}
