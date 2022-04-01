package cli

import (
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"net/url"
	"path/filepath"
)

type ProjectCloneCmd struct {
	Output string `arg:"-o,--output-path" default:"./"`
	Branch string `arg:"-b,--branch" default:"master"`
}

func (b *BitbucketCLI) projectClone(cmd *ProjectCmd) {
	// Clones all the projects
	res, err := b.client.DefaultApi.GetRepositories(cmd.Key)
	if err != nil {
		logrus.Fatal(err)
	}

	repositories, err := bitbucket.GetRepositoriesResponse(res)
	for _, v := range repositories {
		// Get HTTP Clone URL
		var cloneUrl = ""
		for _, cUrl := range v.Links.Clone {
			mUrl, err := url.Parse(cUrl.Href)
			if err != nil {
				continue
			}

			if mUrl.Scheme == "https" {
				cloneUrl = mUrl.String()
				break
			}
		}

		repoPath := filepath.Join(cmd.Clone.Output, v.Slug)
		if err != nil {
			logrus.Warnf("Skipping %s (%s): unable to resolve path: %v", v.Name, v.Slug, err)
			continue
		}

		// Clone repo
		repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			SingleBranch:  true,
			ReferenceName: plumbing.NewBranchReferenceName(cmd.Clone.Branch),
			URL:           cloneUrl,
			Auth:          &b.cloneCredentials,
		})

		if err != nil {
			logrus.Warnf("Unable to clone %s (%s): %v", v.Name, v.Slug, err)
			continue
		}

		head, err := repo.Head()
		if err != nil {
			logrus.Warnf("Unable to get head: %s", err)
			continue
		}

		fmt.Printf("head: %s\n", head.String())
	}
}
