package test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"os"
)

func GetCLI() cli.BitbucketCLI {
	return cli.NewCLI(
		os.Getenv("BITBUCKET_USERNAME"),
		os.Getenv("BITBUCKET_PASSWORD"),
		os.Getenv("BITBUCKET_URL"),
	)
}