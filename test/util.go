package test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"os"
)

func MustGetCLI() cli.BitbucketCLI {
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
