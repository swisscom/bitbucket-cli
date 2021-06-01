package cli

import (
	"fmt"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/sirupsen/logrus"
	"net/url"
)

type ProjectListCmd struct {

}

func (b *BitbucketCLI) projectList(cmd *ProjectCmd){
	if cmd == nil {
		return
	}

	// List project repositories
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
		fmt.Printf("%s\t%s\n", v.Name, cloneUrl)
	}
}