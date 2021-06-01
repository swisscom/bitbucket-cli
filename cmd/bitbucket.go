package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

type ListReposCmd struct {
}

type CloneReposCmd struct {
	Output    string `arg:"-o,--output-path" default:"./"`
	Reference string `arg:"-r,--reference" default:"master"`
}

type ProjectCmd struct {
	Key   string         `arg:"-k,required"`
	List  *ListReposCmd  `arg:"subcommand:list"`
	Clone *CloneReposCmd `arg:"subcommand:clone"`
}

type Args struct {
	Username string      `arg:"-u,--username,env:BITBUCKET_USERNAME"`
	Password string      `arg:"-p,--password,env:BITBUCKET_PASSWORD"`
	Url      string      `arg:"-u,--url,required,env:BITBUCKET_URL"`
	Project  *ProjectCmd `arg:"subcommand:project"`
}

var args Args

func getClient(ctx context.Context, args *Args) *bitbucket.APIClient {
	return bitbucket.NewAPIClient(ctx, bitbucket.NewConfiguration(args.Url))
}

func main() {
	arg.MustParse(&args)
	basicAuth := bitbucket.BasicAuth{UserName: args.Username, Password: args.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Millisecond)
	defer cancel()

	ctx = context.WithValue(ctx, bitbucket.ContextBasicAuth, basicAuth)

	if args.Project != nil {
		if args.Project.Key == "" {
			logrus.Fatal("A project key must be provided")
		}

		if args.Project.List != nil {
			// List project repositories
			c := getClient(ctx, &args)
			res, err := c.DefaultApi.GetRepositories(args.Project.Key)
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
				fmt.Printf("%s\t%s\n", v.Name, cloneUrl)
			}
		} else if args.Project.Clone != nil {
			// Clones all of the projects
			c := getClient(ctx, &args)
			res, err := c.DefaultApi.GetRepositories(args.Project.Key)
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

				repoPath := filepath.Join(args.Project.Clone.Output, v.Slug)
				if err != nil {
					logrus.Warnf("Skipping %s (%s): unable to resolve path: %v", v.Name, v.Slug, err)
					continue
				}

				// Clone repo
				repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{SingleBranch: true,
					URL: cloneUrl, Auth: &http.BasicAuth{Username: args.Username,
						Password: args.Password}})
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
	}
}
