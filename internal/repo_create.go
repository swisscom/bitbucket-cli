package cli

import (
	bitbucket "github.com/gfleury/go-bitbucket-v1"
)

type RepoCreateCmd struct {
	DisplayName   string `arg:"--display-name,required" help:"Display name of the new repository; Bitbucket derives the slug from it"`
	Description   string `arg:"-d,--description" help:"Description of the repository"`
	Forkable      *bool  `arg:"--forkable" help:"Whether the repository can be forked (use --forkable=false to disable)"`
	Public        *bool  `arg:"--public" help:"Whether the repository is publicly accessible (use --public=false to disable)"`
	DefaultBranch string `arg:"--default-branch" help:"Default branch of the new repository, e.g: main"`
	Output        string `arg:"-o,--output" help:"Output format: text (default) or json"`
}

// createBody builds the request body for creating a repository.  Only the
// options the user set are included, so that the server defaults apply to
// everything else.
func createBody(cmd *RepoCreateCmd) map[string]interface{} {
	body := map[string]interface{}{
		"name":  cmd.DisplayName,
		"scmId": "git",
	}
	if cmd.Description != "" {
		body["description"] = cmd.Description
	}
	if cmd.Forkable != nil {
		body["forkable"] = *cmd.Forkable
	}
	if cmd.Public != nil {
		body["public"] = *cmd.Public
	}
	if cmd.DefaultBranch != "" {
		body["defaultBranch"] = cmd.DefaultBranch
	}
	return body
}

// createRepository creates a repository in the project.  See
// parseRepositoryResponse for the two return values.
func (b *BitbucketCLI) createRepository(projectKey string, body map[string]interface{}) (bitbucket.Repository, map[string]interface{}, error) {
	res, err := b.client.DefaultApi.CreateRepositoryWithOptions(projectKey, body, []string{"application/json"})
	if err != nil {
		return bitbucket.Repository{}, nil, err
	}
	return parseRepositoryResponse(res)
}

func (b *BitbucketCLI) repoCreate(cmd *RepoCmd) {
	if cmd == nil || cmd.Create == nil {
		return
	}
	create := cmd.Create

	format, err := outputFormat(create.Output)
	if err != nil {
		b.logger.Fatal(err)
	}

	// The parent's --name (a slug) is deliberately unused here: the new
	// repository is named by --display-name and Bitbucket derives its slug.
	repo, raw, err := b.createRepository(cmd.ProjectKey, createBody(create))
	if err != nil {
		b.logger.Fatalf("unable to create repository \"%s\" in project %s: %v", create.DisplayName, cmd.ProjectKey, err)
	}

	if err := printRepositoryAs(format, repo, raw); err != nil {
		b.logger.Fatalf("unable to print repository: %v", err)
	}
}
