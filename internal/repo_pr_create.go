package cli

import (
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

type RepoPrCreateCmd struct {
	Title       string `arg:"-t,--title,required" help:"Title of this PR"`
	Description string `arg:"-d,--description" help:"Description of the PR"`

	FromRef string `arg:"-F,--from-ref,required" help:"Reference of the incoming PR, e.g: refs/heads/feature-ABC-123"` // e.g: refs/heads/feature-ABC-123
	ToRef   string `arg:"-T,--to-ref,required" help:"Target reference, e.g: refs/heads/master"`

	// From which repo? Defaults to self
	FromRepoKey  string `arg:"-K,--from-key" help:"Project Key of the \"from\" repository"`
	FromRepoSlug string `arg:"-S,--from-slug" help:"Repository slug of the \"from\" repository"`
}

func (b *BitbucketCLI) repoPrCreate(cmd *RepoCmd) {
	if cmd == nil || cmd.PrCmd == nil || cmd.PrCmd.Create == nil {
		return
	}
	create := cmd.PrCmd.Create

	if create.FromRepoKey == "" && create.FromRepoSlug == "" {
		// From = To
		create.FromRepoKey = cmd.ProjectKey
		create.FromRepoSlug = cmd.Slug
	}

	pr := bitbucketv1.PullRequest{
		Title:       create.Title,
		Description: create.Description,
		FromRef: bitbucketv1.PullRequestRef{
			ID: create.FromRef,
			Repository: bitbucketv1.Repository{
				Slug:    create.FromRepoSlug,
				Project: &bitbucketv1.Project{Key: create.FromRepoKey},
			},
		},
		ToRef: bitbucketv1.PullRequestRef{
			ID: create.ToRef,
			Repository: bitbucketv1.Repository{
				Slug:    cmd.Slug,
				Project: &bitbucketv1.Project{Key: cmd.ProjectKey},
			},
		},
	}

	resp, err := b.client.DefaultApi.CreatePullRequest(
		cmd.ProjectKey,
		cmd.Slug,
		pr,
	)
	if err != nil {
		b.logger.Fatalf("unable to create PR: %v", err)
	}

	// Parse resp
	prRes, err := bitbucketv1.GetPullRequestResponse(resp)
	if err != nil {
		b.logger.Fatalf("unable to parse PR: %v", err)
	}

	fmt.Printf("%s", prRes.Links.Self[0].Href)
}
