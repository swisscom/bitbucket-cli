package cli

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// RepoContext holds repository information inferred from the local git repo.
type RepoContext struct {
	ProjectKey string
	Slug       string
	Branch     string // short branch name, e.g. "feature/my-thing"
	BaseURL    string // Bitbucket REST base URL, e.g. "https://host/rest" (HTTPS remotes only)
}

// GetRepoContext opens the git repo at dir (or the nearest parent with a .git),
// reads the current branch and the "origin" remote URL, and parses them into
// a RepoContext.
func GetRepoContext(dir string) (*RepoContext, error) {
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("not inside a git repository: %v", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("could not read HEAD: %v", err)
	}

	branch := ""
	if head.Name().IsBranch() {
		branch = head.Name().Short()
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("no remote named \"origin\": %v", err)
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return nil, fmt.Errorf("remote \"origin\" has no URLs")
	}

	projectKey, slug, baseURL, err := ParseBitbucketRemote(urls[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse Bitbucket project/repo from remote URL %q: %v", urls[0], err)
	}

	return &RepoContext{
		ProjectKey: projectKey,
		Slug:       slug,
		Branch:     branch,
		BaseURL:    baseURL,
	}, nil
}

// ParseBitbucketRemote parses a Bitbucket Server remote URL and returns the
// project key, repository slug, and (for HTTPS remotes) the REST API base URL.
//
// Supported formats:
//
//	HTTPS (root deployment):    https://host/scm/PROJECT/repo.git
//	HTTPS (context path):       https://host/context/scm/PROJECT/repo.git
//	SSH URL with port:          ssh://git@host:7999/PROJECT/repo.git
//	SSH URL without port:       ssh://git@host/PROJECT/repo.git
//	git+ssh URL:                git+ssh://git@host:7999/PROJECT/repo.git
//	SCP-style:                  git@host:PROJECT/repo.git
func ParseBitbucketRemote(remoteURL string) (projectKey, slug, baseURL string, err error) {
	switch {
	case strings.HasPrefix(remoteURL, "https://"), strings.HasPrefix(remoteURL, "http://"):
		return parseHTTPSRemote(remoteURL)
	case strings.HasPrefix(remoteURL, "ssh://"), strings.HasPrefix(remoteURL, "git+ssh://"):
		return parseSSHURLRemote(remoteURL)
	default:
		// Fall through to SCP-style: [user@]host:path
		return parseSCPRemote(remoteURL)
	}
}

// parseHTTPSRemote handles http(s):// URLs.
// Bitbucket Server HTTPS clone URLs always contain "/scm/" as the separator
// between an optional context path and the project/repo path.
func parseHTTPSRemote(remoteURL string) (projectKey, slug, baseURL string, err error) {
	u, err := url.Parse(remoteURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid URL: %v", err)
	}

	// Split on "/scm/" to separate the optional context path from the project/repo.
	const scmSegment = "/scm/"
	idx := strings.Index(u.Path, scmSegment)
	if idx < 0 {
		return "", "", "", fmt.Errorf("URL does not contain /scm/ — is this a Bitbucket Server clone URL?")
	}

	contextPath := u.Path[:idx]            // e.g. "" or "/bitbucket"
	repoPath := u.Path[idx+len(scmSegment):] // e.g. "PROJECT/repo.git"

	projectKey, slug, err = splitProjectSlug(repoPath)
	if err != nil {
		return "", "", "", err
	}

	// Strip userinfo (username in clone URLs) from host for the base URL.
	baseURL = fmt.Sprintf("%s://%s%s/rest", u.Scheme, u.Hostname(), contextPath)
	if u.Port() != "" {
		baseURL = fmt.Sprintf("%s://%s:%s%s/rest", u.Scheme, u.Hostname(), u.Port(), contextPath)
	}

	return projectKey, slug, baseURL, nil
}

// parseSSHURLRemote handles ssh:// and git+ssh:// URLs.
// Bitbucket Server SSH URLs carry the project/repo directly in the path with
// no context path component.
func parseSSHURLRemote(remoteURL string) (projectKey, slug, baseURL string, err error) {
	// Normalise git+ssh:// → ssh:// so url.Parse works cleanly.
	normalized := strings.TrimPrefix(remoteURL, "git+")

	u, err := url.Parse(normalized)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid SSH URL: %v", err)
	}

	// Path is "/PROJECT/repo[.git]"; strip the leading slash.
	projectKey, slug, err = splitProjectSlug(strings.TrimPrefix(u.Path, "/"))
	if err != nil {
		return "", "", "", err
	}

	return projectKey, slug, "", nil // base URL not determinable from SSH
}

// parseSCPRemote handles SCP-style URLs: [user@]host:PROJECT/repo[.git]
func parseSCPRemote(remoteURL string) (projectKey, slug, baseURL string, err error) {
	colonIdx := strings.Index(remoteURL, ":")
	if colonIdx < 0 {
		return "", "", "", fmt.Errorf("not a recognised git remote URL: %q", remoteURL)
	}

	path := remoteURL[colonIdx+1:]
	projectKey, slug, err = splitProjectSlug(path)
	if err != nil {
		return "", "", "", err
	}

	return projectKey, slug, "", nil // base URL not determinable from SCP
}

// splitProjectSlug splits a "PROJECT/repo[.git]" string into its two parts.
func splitProjectSlug(path string) (projectKey, slug string, err error) {
	path = strings.TrimSuffix(path, ".git")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected PROJECT/repo in %q", path)
	}
	return parts[0], parts[1], nil
}

// GetSingleCommitMessage looks at the commits in fromRef that are not
// reachable from toRef.  If there is exactly one such commit it returns its
// subject line and body so callers can pre-populate PR title/description.
// Returns found=false (without error) when the count is not exactly one.
func GetSingleCommitMessage(repoPath, fromRef, toRef string) (subject, body string, found bool, err error) {
	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", "", false, fmt.Errorf("not inside a git repository: %v", err)
	}

	fromHash, err := repo.ResolveRevision(plumbing.Revision(fromRef))
	if err != nil {
		return "", "", false, fmt.Errorf("could not resolve %q: %v", fromRef, err)
	}

	toHash, err := repo.ResolveRevision(plumbing.Revision(toRef))
	if err != nil {
		return "", "", false, fmt.Errorf("could not resolve %q: %v", toRef, err)
	}

	// Build the set of commits reachable from toRef so we can stop the walk.
	toAncestors := make(map[plumbing.Hash]bool)
	toIter, err := repo.Log(&git.LogOptions{From: *toHash})
	if err != nil {
		return "", "", false, err
	}
	_ = toIter.ForEach(func(c *object.Commit) error {
		toAncestors[c.Hash] = true
		return nil
	})

	// Walk from fromRef, collecting commits not in toAncestors.
	// Stop as soon as we have seen more than one (no need to count further).
	var unique []*object.Commit
	fromIter, err := repo.Log(&git.LogOptions{From: *fromHash})
	if err != nil {
		return "", "", false, err
	}
	_ = fromIter.ForEach(func(c *object.Commit) error {
		if toAncestors[c.Hash] {
			return storer.ErrStop
		}
		unique = append(unique, c)
		if len(unique) > 1 {
			return storer.ErrStop
		}
		return nil
	})

	if len(unique) != 1 {
		return "", "", false, nil
	}

	subject, body = splitCommitMessage(unique[0].Message)
	return subject, body, true, nil
}

// splitCommitMessage splits a raw commit message into subject and body,
// following the conventional blank-line separator.
func splitCommitMessage(msg string) (subject, body string) {
	msg = strings.TrimRight(msg, "\n")
	parts := strings.SplitN(msg, "\n", 2)
	subject = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}
	return subject, body
}
