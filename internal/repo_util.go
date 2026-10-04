package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	bitbucket "github.com/gfleury/go-bitbucket-v1"
)

const (
	outputText = "text"
	outputJson = "json"
)

// outputFormat validates the value of an --output flag.  An empty value
// selects the default text output.
func outputFormat(value string) (string, error) {
	if value == "" {
		return outputText, nil
	}

	validFormats := []string{outputText, outputJson}
	format := strings.ToLower(value)
	if !validValue(format, validFormats) {
		return "", fmt.Errorf("invalid value \"%s\" for output: accepted values are: %s",
			value,
			strings.Join(validFormats, ", "),
		)
	}
	return format, nil
}

// parseRepositoryResponse decodes a single-repository response into the
// library's typed struct and also returns the raw response map, which
// carries fields the struct lacks (such as the description).
func parseRepositoryResponse(res *bitbucket.APIResponse) (bitbucket.Repository, map[string]interface{}, error) {
	repo, err := bitbucket.GetRepositoryResponse(res)
	if err != nil {
		return bitbucket.Repository{}, nil, fmt.Errorf("unable to parse repository response: %v", err)
	}
	return repo, res.Values, nil
}

// repositoryDescription reads the description from the raw response, since
// the library's Repository struct has no field for it.
func repositoryDescription(raw map[string]interface{}) string {
	description, _ := raw["description"].(string)
	return description
}

// cloneUrl returns the clone link of the repository whose URL uses the given
// scheme (e.g. "https" or "ssh"), or "" if there is none.
func cloneUrl(repo bitbucket.Repository, scheme string) string {
	if repo.Links == nil {
		return ""
	}

	for _, link := range repo.Links.Clone {
		mUrl, err := url.Parse(link.Href)
		if err != nil {
			continue
		}
		if mUrl.Scheme == scheme {
			return mUrl.String()
		}
	}
	return ""
}

// selfUrl returns the web link of the repository, or "" if there is none.
func selfUrl(repo bitbucket.Repository) string {
	if repo.Links == nil || len(repo.Links.Self) == 0 {
		return ""
	}
	return repo.Links.Self[0].Href
}

// printRepository writes a human readable summary of the repository to
// stdout.
func printRepository(repo bitbucket.Repository, raw map[string]interface{}) {
	projectKey := ""
	if repo.Project != nil {
		projectKey = repo.Project.Key
	}

	lines := []struct {
		label string
		value string
	}{
		{"Name", repo.Name},
		{"Slug", repo.Slug},
		{"ID", fmt.Sprintf("%d", repo.ID)},
		{"Project", projectKey},
		{"State", repo.State},
		{"Public", fmt.Sprintf("%t", repo.Public)},
		{"Forkable", fmt.Sprintf("%t", repo.Forkable)},
		{"Description", repositoryDescription(raw)},
		{"Clone (https)", cloneUrl(repo, "https")},
		{"Clone (ssh)", cloneUrl(repo, "ssh")},
		{"Web", selfUrl(repo)},
	}

	for _, line := range lines {
		fmt.Printf("%-14s %s\n", line.label+":", line.value)
	}
}

// printJson writes v to stdout as indented JSON followed by a newline.
func printJson(v interface{}) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", out)
	return err
}

// printRepositoryAs prints the repository in the requested output format.
func printRepositoryAs(format string, repo bitbucket.Repository, raw map[string]interface{}) error {
	if format == outputJson {
		return printJson(raw)
	}
	printRepository(repo, raw)
	return nil
}
