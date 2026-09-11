package cli

import (
	"fmt"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"strings"
)

type RepoPrCreateCmd struct {
	Title       string `arg:"-t,--title" help:"Title of this PR; defaults to the commit subject when the branch has exactly one commit"`
	Description string `arg:"-d,--description" help:"Description of the PR; defaults to the commit body when the branch has exactly one commit"`

	FromRef string `arg:"-F,--from-ref" help:"Reference of the incoming PR, e.g: refs/heads/feature-ABC-123; defaults to the current branch"` // e.g: refs/heads/feature-ABC-123
	ToRef   string `arg:"-T,--to-ref,required" help:"Target reference, e.g: refs/heads/master"`

	// From which repo? Defaults to self
	FromRepoKey  string `arg:"-K,--from-key" help:"Project key of the \"from\" repository"`
	FromRepoSlug string `arg:"-S,--from-slug" help:"Repository slug of the \"from\" repository"`

	Reviewers string `arg:"-r,--reviewers,env:BITBUCKET_REVIEWERS" help:"Comma separated list of reviewers"`
}

func (b BitbucketCLI) GetReviewers(revList string) []bitbucketv1.UserWithMetadata {
	if revList == "" {
		return nil
	}
	var reviewers []bitbucketv1.UserWithMetadata
	for _, user := range strings.Split(revList, ",") {
		if usersResponse, err := b.client.DefaultApi.GetUsers(map[string]interface{}{"filter": user}); err != nil {
			b.logger.Fatalf("Error while retrieving user %s: %e", user, err)
		} else if users, err := bitbucketv1.GetUsersResponse(usersResponse); err != nil {
			b.logger.Fatalf("Error while parsing list of users for user %s: %e", user, err)
		} else if len(users) == 0 {
			b.logger.Fatalf("user %s does not exist", user)
		} else if len(users) > 1 {
			var found []string
			for _, bbUser := range users {
				found = append(found, fmt.Sprintf("%s: %s (%s)", bbUser.Slug, bbUser.Name, bbUser.EmailAddress))
			}
			b.logger.Fatalf("multiple users found for user %s: %s", user, strings.Join(found, ", "))
		} else {
			bbUser := users[0]
			reviewers = append(reviewers, bitbucketv1.UserWithMetadata{
				User: bitbucketv1.UserWithLinks{
					Name:         bbUser.Name,
					EmailAddress: bbUser.EmailAddress,
					Slug:         bbUser.Slug,
				},
			})
		}
	}
	return reviewers
}

func (b *BitbucketCLI) repoPrCreate(cmd *RepoCmd) {
	if cmd == nil || cmd.PrCmd == nil || cmd.PrCmd.Create == nil {
		return
	}
	create := cmd.PrCmd.Create

	if create.FromRef == "" {
		ctx, err := GetRepoContext(".")
		if err != nil {
			b.logger.Fatalf("--from-ref not specified and could not detect current branch: %v", err)
		}
		if ctx.Branch == "" {
			b.logger.Fatal("--from-ref not specified and HEAD is not on a branch (detached HEAD).")
		}
		create.FromRef = "refs/heads/" + ctx.Branch
	}

	if create.Title == "" {
		subject, body, found, err := GetSingleCommitMessage(".", create.FromRef, create.ToRef)
		if err != nil {
			b.logger.Fatalf("--title not specified and could not read commit history: %v", err)
		}
		if !found {
			b.logger.Fatal("--title not specified and branch does not have exactly one commit ahead of the target.")
		}
		create.Title = subject
		if create.Description == "" {
			create.Description = body
		}
	}

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
		Reviewers: b.GetReviewers(create.Reviewers),
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
