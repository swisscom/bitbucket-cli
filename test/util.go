package test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"os"
	"testing"
)

func SkipIfNoServer(t *testing.T) {
	t.Helper()
	if os.Getenv("BITBUCKET_URL") == "" {
		t.Skip("BITBUCKET_URL not set, skipping integration test")
	}
}

func MustGetCLI() *cli.BitbucketCLI {
	c, err := cli.NewCLI(
		&cli.BasicAuth{
			Username: os.Getenv("BITBUCKET_USERNAME"),
			Password: os.Getenv("BITBUCKET_PASSWORD"),
		},
		os.Getenv("BITBUCKET_URL"),
	)
	if err != nil {
		panic(err)
	}
	return c
}
