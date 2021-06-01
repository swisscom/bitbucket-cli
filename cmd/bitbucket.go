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
}

var args Args

func main() {
	arg.MustParse(&args)
	logger := logrus.New()

	c := cli.NewCLI(args.Username, args.Password, args.Url)
	c.SetLogger(logger)

	if args.Project != nil {
		c.RunProjectCmd(args.Project)
	}
}
