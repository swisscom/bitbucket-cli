package cli

import (
	"context"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/sirupsen/logrus"
	"time"
)

type BitbucketCLI struct {
	username string
	password string
	repoUrl  string

	client *bitbucket.APIClient
	logger *logrus.Logger
}

func (b BitbucketCLI) SetLogger(logger *logrus.Logger) {
	if logger == nil {
		// We don't set nil loggers
		return
	}
	b.logger = logger
}

func NewCLI(username string, password string, repoUrl string) BitbucketCLI {
	basicAuth := bitbucket.BasicAuth{UserName: username, Password: password}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ctx = context.WithValue(ctx, bitbucket.ContextBasicAuth, basicAuth)
	c := bitbucket.NewAPIClient(ctx, bitbucket.NewConfiguration(repoUrl))
	logger := logrus.New()

	return BitbucketCLI{
		username: username,
		password: password,
		repoUrl:  repoUrl,
		client:   c,
		logger:   logger,
	}
}
