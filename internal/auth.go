package cli

import (
	"context"
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"net/http"
)

type BasicAuth struct {
	Username string
	Password string
}

func (b BasicAuth) AddHeaders(req *http.Request) {
	req.SetBasicAuth(b.Username, b.Password)
}

func (b BasicAuth) GetCloneCredentials() git_http.BasicAuth {
	return git_http.BasicAuth{Username: b.Username, Password: b.Password}
}

func (b BasicAuth) GetContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, bitbucket.ContextBasicAuth,
		bitbucket.BasicAuth{
			UserName: b.Username,
			Password: b.Password,
		})
}

type AccessToken struct {
	Username    string
	AccessToken string
}

func (a AccessToken) AddHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken))
}

func (a AccessToken) GetCloneCredentials() git_http.BasicAuth {
	return git_http.BasicAuth{
		Username: a.Username,
		Password: a.AccessToken,
	}
}

func (a AccessToken) GetContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, bitbucket.ContextAccessToken, a.AccessToken)
}

var _ Authenticator = (*BasicAuth)(nil)
var _ Authenticator = (*AccessToken)(nil)

type Authenticator interface {
	GetContext(ctx context.Context) context.Context
	GetCloneCredentials() git_http.BasicAuth
	AddHeaders(req *http.Request)
}
