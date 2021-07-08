package main

import (
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	cli "github.com/swisscom/bitbucket-cli/internal"
)

type Args struct {
	Username string          `arg:"-u,--username,env:BITBUCKET_USERNAME"`
	Password string          `arg:"-p,--password,env:BITBUCKET_PASSWORD"`
	Url      string          `arg:"-u,--url,required,env:BITBUCKET_URL"`
	Project  *cli.ProjectCmd `arg:"subcommand:project"`
	Repo     *cli.RepoCmd    `arg:"subcommand:repo"`
}

var args Args

func main() {
	p := arg.MustParse(&args)
	logger := logrus.New()

	c := cli.NewCLI(args.Username, args.Password, args.Url)
	c.SetLogger(logger)

	if args.Project != nil {
		c.RunProjectCmd(args.Project)
		return
	}

	if args.Repo != nil {
		c.RunRepoCmd(args.Repo)
		return
	}

	p.Fail("Command must be specified")

}
