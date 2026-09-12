package cli

import (
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

type PrCreateCmd struct {
	ProjectKey  string `arg:"-k,--key,required,env:BITBUCKET_PROJECT" help:"Project key"`
	Slug        string `arg:"-n,--name,required,env:BITBUCKET_REPO" help:"Repository slug"`
	Title       string `arg:"-t,--title,required" help:"Title of this PR"`
	Description string `arg:"-d,--description" help:"Description of the PR"`

	FromRef string `arg:"-F,--from-ref,required" help:"Reference of the incoming PR, e.g: refs/heads/feature-ABC-123"`
	ToRef   string `arg:"-T,--to-ref,required" help:"Target reference, e.g: refs/heads/master"`

	// From which repo? Defaults to self
	FromRepoKey  string `arg:"-K,--from-key" help:"Project Key of the \"from\" repository"`
	FromRepoSlug string `arg:"-S,--from-slug" help:"Repository slug of the \"from\" repository"`

	Reviewers string `arg:"-r,--reviewers,env:BITBUCKET_REVIEWERS" help:"Comma separated list of reviewers"`
}

func (b *BitbucketCLI) RunPRCreateCmd(cmd *PrCreateCmd) {
	if cmd == nil {
		return
	}

	if cmd.FromRepoKey == "" && cmd.FromRepoSlug == "" {
		// From = To
		cmd.FromRepoKey = cmd.ProjectKey
		cmd.FromRepoSlug = cmd.Slug
	}

	pr := bitbucketv1.PullRequest{
		Title:       cmd.Title,
		Description: cmd.Description,
		FromRef: bitbucketv1.PullRequestRef{
			ID: cmd.FromRef,
			Repository: bitbucketv1.Repository{
				Slug:    cmd.FromRepoSlug,
				Project: &bitbucketv1.Project{Key: cmd.FromRepoKey},
			},
		},
		ToRef: bitbucketv1.PullRequestRef{
			ID: cmd.ToRef,
			Repository: bitbucketv1.Repository{
				Slug:    cmd.Slug,
				Project: &bitbucketv1.Project{Key: cmd.ProjectKey},
			},
		},
		Reviewers: b.GetReviewers(cmd.Reviewers),
	}

	resp, err := b.client.DefaultApi.CreatePullRequest(cmd.ProjectKey, cmd.Slug, pr)
	if err != nil {
		b.logger.Fatalf("unable to create PR: %v", err)
	}

	prRes, err := bitbucketv1.GetPullRequestResponse(resp)
	if err != nil {
		b.logger.Fatalf("unable to parse PR: %v", err)
	}

	fmt.Printf("%s", prRes.Links.Self[0].Href)
}
