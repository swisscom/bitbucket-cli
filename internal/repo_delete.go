package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RepoDeleteCmd struct {
	Yes bool `arg:"-y,--yes" help:"Skip the confirmation prompt"`
}

// maxErrorBodySize bounds how much of an error response body is echoed back
// to the user.
const maxErrorBodySize = 64 * 1024

var errDeletionAborted = errors.New("deletion aborted: confirmation did not match; pass --yes to skip the prompt")

// confirmDeletion asks the user, via out, to type the repository slug and
// reads the answer from in.  It returns an error unless the answer matches.
func confirmDeletion(in io.Reader, out io.Writer, projectKey string, slug string, displayName string) error {
	_, err := fmt.Fprintf(out, "Delete repository %s/%s (\"%s\")? Type the slug to confirm: ", projectKey, slug, displayName)
	if err != nil {
		return fmt.Errorf("unable to write confirmation prompt: %v", err)
	}

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("unable to read confirmation: %v", err)
	}

	if strings.TrimSpace(line) != slug {
		return errDeletionAborted
	}
	return nil
}

// deleteRepository schedules the repository for deletion.
//
// This deliberately bypasses the library's DeleteRepository: Bitbucket
// answers with an empty body, which the library cannot decode, and the two
// meaningful status codes (202 scheduled, 204 does not exist) would be
// indistinguishable.
func (b *BitbucketCLI) deleteRepository(projectKey string, slug string) error {
	deleteUrl := fmt.Sprintf("%s/api/1.0/projects/%s/repos/%s",
		b.apiBaseUrl(),
		url.PathEscape(projectKey),
		url.PathEscape(slug),
	)
	b.logger.Debugf("DELETE %s", deleteUrl)

	req, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}

	res, err := b.doReq(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %v", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusAccepted:
		return nil
	case http.StatusNoContent:
		return fmt.Errorf("repository %s/%s does not exist", projectKey, slug)
	default:
		body, readErr := io.ReadAll(io.LimitReader(res.Body, maxErrorBodySize))
		if readErr != nil {
			return fmt.Errorf("unexpected status %s (unable to read response body: %v)", res.Status, readErr)
		}
		return fmt.Errorf("unexpected status %s: %s", res.Status, strings.TrimSpace(string(body)))
	}
}

func (b *BitbucketCLI) repoDelete(cmd *RepoCmd) {
	if cmd == nil || cmd.Delete == nil {
		return
	}

	if !cmd.Delete.Yes {
		// Fetch the repository first: it proves the repository exists and
		// gives us the display name for the prompt.
		repo, _, err := b.getRepository(cmd.ProjectKey, cmd.Slug)
		if err != nil {
			b.logger.Fatalf("unable to get repository %s/%s: %v", cmd.ProjectKey, cmd.Slug, err)
		}

		// The prompt goes to stderr so that stdout stays clean for scripts.
		if err := confirmDeletion(os.Stdin, os.Stderr, cmd.ProjectKey, cmd.Slug, repo.Name); err != nil {
			b.logger.Fatal(err)
		}
	}

	if err := b.deleteRepository(cmd.ProjectKey, cmd.Slug); err != nil {
		b.logger.Fatalf("unable to delete repository %s/%s: %v", cmd.ProjectKey, cmd.Slug, err)
	}

	fmt.Printf("repository %s/%s scheduled for deletion\n", cmd.ProjectKey, cmd.Slug)
}
