package cli

import (
	"errors"

	bitbucket "github.com/gfleury/go-bitbucket-v1"
)

type RepoUpdateCmd struct {
	NewName       string  `arg:"--new-name" help:"New name of the repository (renaming may change the slug)"`
	Description   *string `arg:"-d,--description" help:"New description of the repository (pass \"\" to clear it)"`
	Forkable      *bool   `arg:"--forkable" help:"Whether the repository can be forked (use --forkable=false to disable)"`
	Public        *bool   `arg:"--public" help:"Whether the repository is publicly accessible (use --public=false to disable)"`
	DefaultBranch string  `arg:"--default-branch" help:"New default branch, e.g: main"`
	ToProject     string  `arg:"--to-project" help:"Key of the project to move the repository to"`
	Output        string  `arg:"-o,--output" help:"Output format: text (default) or json"`
}

var errNothingToUpdate = errors.New(
	"nothing to update: specify at least one of --new-name, --description, --forkable, --public, --default-branch or --to-project",
)

// updateBody builds the request body for updating a repository.  Only the
// options the user set are included, so that everything else keeps its
// current value.
func updateBody(cmd *RepoUpdateCmd) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if cmd.NewName != "" {
		body["name"] = cmd.NewName
	}
	if cmd.Description != nil {
		body["description"] = *cmd.Description
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
	if cmd.ToProject != "" {
		body["project"] = map[string]interface{}{"key": cmd.ToProject}
	}

	if len(body) == 0 {
		return nil, errNothingToUpdate
	}
	return body, nil
}

// updateRepository updates a repository.  See parseRepositoryResponse for
// the two return values.
func (b *BitbucketCLI) updateRepository(projectKey string, slug string, body map[string]interface{}) (bitbucket.Repository, map[string]interface{}, error) {
	res, err := b.client.DefaultApi.UpdateRepositoryWithOptions(projectKey, slug, body, []string{"application/json"})
	if err != nil {
		return bitbucket.Repository{}, nil, err
	}
	return parseRepositoryResponse(res)
}

func (b *BitbucketCLI) repoUpdate(cmd *RepoCmd) {
	if cmd == nil || cmd.Update == nil {
		return
	}
	update := cmd.Update

	format, err := outputFormat(update.Output)
	if err != nil {
		b.logger.Fatal(err)
	}

	body, err := updateBody(update)
	if err != nil {
		b.logger.Fatal(err)
	}

	repo, raw, err := b.updateRepository(cmd.ProjectKey, cmd.Slug, body)
	if err != nil {
		b.logger.Fatalf("unable to update repository %s/%s: %v", cmd.ProjectKey, cmd.Slug, err)
	}

	if err := printRepositoryAs(format, repo, raw); err != nil {
		b.logger.Fatalf("unable to print repository: %v", err)
	}
}
