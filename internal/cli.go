package cli

import (
	"context"
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// requestTimeout bounds every single HTTP request made to Bitbucket.
const requestTimeout = 10 * time.Second

type BitbucketCLI struct {
	cloneCredentials gitHttp.BasicAuth
	restUrl          *url.URL
	client           *bitbucket.APIClient
	logger           *logrus.Logger
	httpClient       *http.Client
	auth             Authenticator
}

func (b *BitbucketCLI) SetLogger(logger *logrus.Logger) {
	if logger == nil {
		// We don't set nil loggers
		return
	}
	b.logger = logger
}

func NewCLI(auth Authenticator, restUrl string) (*BitbucketCLI, error) {
	mUrl, err := url.Parse(restUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %v", err)
	}

	// The timeout lives on the HTTP client so that it applies per request.
	// A deadline on the context would start counting at construction time
	// and expire for every request made after it.
	httpClient := &http.Client{Timeout: requestTimeout}
	ctx := auth.GetContext(context.Background())
	c := bitbucket.NewAPIClient(ctx, bitbucket.NewConfiguration(
		strings.TrimRight(mUrl.String(), "/"), // https://git.example.com/rest/ -> https://git.example.com/rest
		func(cfg *bitbucket.Configuration) {
			cfg.HTTPClient = httpClient
		},
	))
	logger := logrus.New()

	return &BitbucketCLI{
		cloneCredentials: auth.GetCloneCredentials(),
		restUrl:          mUrl,
		auth:             auth,
		client:           c,
		logger:           logger,
		httpClient:       httpClient,
	}, nil
}
