package cli

import (
	bitbucket "github.com/gfleury/go-bitbucket-v1"
)

type RepoGetCmd struct {
	Output string `arg:"-o,--output" help:"Output format: text (default) or json"`
}

// getRepository fetches a repository.  See parseRepositoryResponse for the
// two return values.
func (b *BitbucketCLI) getRepository(projectKey string, slug string) (bitbucket.Repository, map[string]interface{}, error) {
	res, err := b.client.DefaultApi.GetRepository(projectKey, slug)
	if err != nil {
		return bitbucket.Repository{}, nil, err
	}
	return parseRepositoryResponse(res)
}

func (b *BitbucketCLI) repoGet(cmd *RepoCmd) {
	if cmd == nil || cmd.Get == nil {
		return
	}

	format, err := outputFormat(cmd.Get.Output)
	if err != nil {
		b.logger.Fatal(err)
	}

	repo, raw, err := b.getRepository(cmd.ProjectKey, cmd.Slug)
	if err != nil {
		b.logger.Fatalf("unable to get repository %s/%s: %v", cmd.ProjectKey, cmd.Slug, err)
	}

	if err := printRepositoryAs(format, repo, raw); err != nil {
		b.logger.Fatalf("unable to print repository: %v", err)
	}
}
