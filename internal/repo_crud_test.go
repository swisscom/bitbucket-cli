package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// These tests exercise the repository commands against a fake Bitbucket
// served by httptest; they need no real server and never skip.

const repositoryJson = `{
  "slug": "my-repo",
  "id": 42,
  "name": "My Repo",
  "description": "A test repository",
  "scmId": "git",
  "state": "AVAILABLE",
  "statusMessage": "Available",
  "forkable": true,
  "project": {"key": "PRJ", "id": 1, "name": "Project", "public": false, "type": "NORMAL"},
  "public": false,
  "links": {
    "clone": [
      {"href": "ssh://git@bitbucket.example.com:7999/prj/my-repo.git", "name": "ssh"},
      {"href": "https://bitbucket.example.com/scm/prj/my-repo.git", "name": "http"}
    ],
    "self": [{"href": "https://bitbucket.example.com/projects/PRJ/repos/my-repo/browse"}]
  }
}`

const errorJson = `{"errors":[{"context":null,"message":"You are not permitted to access this resource","exceptionName":"com.atlassian.bitbucket.AuthorisationException"}]}`

// recordedRequest captures what the fake server received.
type recordedRequest struct {
	method        string
	path          string
	contentType   string
	authorization string
	atlassianTok  string
	body          map[string]interface{}
}

// newTestCLI starts a fake server answering every request with status and
// response, and returns a CLI pointed at it together with the recorder.
func newTestCLI(t *testing.T, status int, response string) (*BitbucketCLI, *recordedRequest) {
	t.Helper()

	recorded := &recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorded.method = r.Method
		recorded.path = r.URL.Path
		recorded.contentType = r.Header.Get("Content-Type")
		recorded.authorization = r.Header.Get("Authorization")
		recorded.atlassianTok = r.Header.Get("X-Atlassian-Token")

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("unable to read request body: %v", err)
		}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &recorded.body); err != nil {
				t.Errorf("request body is not JSON: %v (%s)", err, raw)
			}
		}

		w.WriteHeader(status)
		if _, err := io.WriteString(w, response); err != nil {
			t.Errorf("unable to write response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	c, err := NewCLI(&BasicAuth{Username: "user", Password: "secret"}, server.URL)
	if err != nil {
		t.Fatalf("unable to create CLI: %v", err)
	}
	return c, recorded
}

func boolPtr(v bool) *bool {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func TestCreateBody(t *testing.T) {
	tests := []struct {
		name string
		cmd  RepoCreateCmd
		want map[string]interface{}
	}{
		{
			name: "defaults",
			cmd:  RepoCreateCmd{DisplayName: "My Repo"},
			want: map[string]interface{}{"name": "My Repo", "scmId": "git"},
		},
		{
			name: "explicit false is sent",
			cmd:  RepoCreateCmd{DisplayName: "My Repo", Forkable: boolPtr(false)},
			want: map[string]interface{}{"name": "My Repo", "scmId": "git", "forkable": false},
		},
		{
			name: "all options",
			cmd: RepoCreateCmd{
				DisplayName:   "My Repo",
				Description:   "desc",
				Forkable:      boolPtr(true),
				Public:        boolPtr(true),
				DefaultBranch: "main",
			},
			want: map[string]interface{}{
				"name":          "My Repo",
				"scmId":         "git",
				"description":   "desc",
				"forkable":      true,
				"public":        true,
				"defaultBranch": "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createBody(&tt.cmd)
			assertJsonEqual(t, tt.want, got)
		})
	}
}

func TestUpdateBody(t *testing.T) {
	tests := []struct {
		name string
		cmd  RepoUpdateCmd
		want map[string]interface{}
	}{
		{
			name: "rename only",
			cmd:  RepoUpdateCmd{NewName: "New Name"},
			want: map[string]interface{}{"name": "New Name"},
		},
		{
			name: "clear description",
			cmd:  RepoUpdateCmd{Description: stringPtr("")},
			want: map[string]interface{}{"description": ""},
		},
		{
			name: "explicit false is sent",
			cmd:  RepoUpdateCmd{Forkable: boolPtr(false), Public: boolPtr(false)},
			want: map[string]interface{}{"forkable": false, "public": false},
		},
		{
			name: "move to project",
			cmd:  RepoUpdateCmd{ToProject: "OTHER", DefaultBranch: "main"},
			want: map[string]interface{}{
				"project":       map[string]interface{}{"key": "OTHER"},
				"defaultBranch": "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updateBody(&tt.cmd)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertJsonEqual(t, tt.want, got)
		})
	}

	t.Run("nothing to update", func(t *testing.T) {
		_, err := updateBody(&RepoUpdateCmd{})
		if err != errNothingToUpdate {
			t.Fatalf("expected errNothingToUpdate, got %v", err)
		}
	})
}

func TestOutputFormat(t *testing.T) {
	for input, want := range map[string]string{"": outputText, "text": outputText, "JSON": outputJson} {
		got, err := outputFormat(input)
		if err != nil {
			t.Errorf("outputFormat(%q): unexpected error: %v", input, err)
		}
		if got != want {
			t.Errorf("outputFormat(%q) = %q, want %q", input, got, want)
		}
	}
	if _, err := outputFormat("yaml"); err == nil {
		t.Error("outputFormat(\"yaml\"): expected an error")
	}
}

func TestGetRepository(t *testing.T) {
	c, recorded := newTestCLI(t, http.StatusOK, repositoryJson)

	repo, raw, err := c.getRepository("PRJ", "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertRequest(t, recorded, http.MethodGet, "/api/1.0/projects/PRJ/repos/my-repo")
	if repo.Name != "My Repo" || repo.Slug != "my-repo" || repo.ID != 42 {
		t.Errorf("unexpected repository: %+v", repo)
	}
	if repo.Project == nil || repo.Project.Key != "PRJ" {
		t.Errorf("unexpected project: %+v", repo.Project)
	}
	if got := repositoryDescription(raw); got != "A test repository" {
		t.Errorf("description = %q", got)
	}
	if got := cloneUrl(repo, "https"); got != "https://bitbucket.example.com/scm/prj/my-repo.git" {
		t.Errorf("https clone url = %q", got)
	}
	if got := cloneUrl(repo, "ssh"); got != "ssh://git@bitbucket.example.com:7999/prj/my-repo.git" {
		t.Errorf("ssh clone url = %q", got)
	}
	if got := selfUrl(repo); got != "https://bitbucket.example.com/projects/PRJ/repos/my-repo/browse" {
		t.Errorf("self url = %q", got)
	}
}

func TestGetRepositoryMinimalResponse(t *testing.T) {
	c, _ := newTestCLI(t, http.StatusOK, `{"slug": "bare", "name": "bare"}`)

	repo, raw, err := c.getRepository("PRJ", "bare")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := repositoryDescription(raw); got != "" {
		t.Errorf("description = %q, want empty", got)
	}
	if got := cloneUrl(repo, "https"); got != "" {
		t.Errorf("clone url = %q, want empty", got)
	}
	if got := selfUrl(repo); got != "" {
		t.Errorf("self url = %q, want empty", got)
	}
	// Must not panic on the missing optional fields.
	printRepository(repo, raw)
}

func TestGetRepositoryError(t *testing.T) {
	c, _ := newTestCLI(t, http.StatusNotFound, `{"errors":[{"message":"Repository PRJ/missing does not exist."}]}`)

	_, _, err := c.getRepository("PRJ", "missing")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error does not carry the server message: %v", err)
	}
}

func TestCreateRepository(t *testing.T) {
	c, recorded := newTestCLI(t, http.StatusCreated, repositoryJson)

	body := createBody(&RepoCreateCmd{DisplayName: "My Repo", Forkable: boolPtr(false)})
	repo, _, err := c.createRepository("PRJ", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertRequest(t, recorded, http.MethodPost, "/api/1.0/projects/PRJ/repos")
	assertJsonContentType(t, recorded)
	assertJsonEqual(t, map[string]interface{}{"name": "My Repo", "scmId": "git", "forkable": false}, recorded.body)
	if repo.Slug != "my-repo" {
		t.Errorf("slug = %q", repo.Slug)
	}
}

func TestCreateRepositoryError(t *testing.T) {
	c, _ := newTestCLI(t, http.StatusUnauthorized, errorJson)

	_, _, err := c.createRepository("PRJ", createBody(&RepoCreateCmd{DisplayName: "My Repo"}))
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "not permitted") {
		t.Errorf("error does not carry the server message: %v", err)
	}
}

func TestUpdateRepository(t *testing.T) {
	c, recorded := newTestCLI(t, http.StatusCreated, repositoryJson)

	body, err := updateBody(&RepoUpdateCmd{Description: stringPtr(""), ToProject: "OTHER"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	repo, _, err := c.updateRepository("PRJ", "my-repo", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertRequest(t, recorded, http.MethodPut, "/api/1.0/projects/PRJ/repos/my-repo")
	assertJsonContentType(t, recorded)
	assertJsonEqual(t, map[string]interface{}{
		"description": "",
		"project":     map[string]interface{}{"key": "OTHER"},
	}, recorded.body)
	if repo.Name != "My Repo" {
		t.Errorf("name = %q", repo.Name)
	}
}

func TestDeleteRepository(t *testing.T) {
	t.Run("scheduled", func(t *testing.T) {
		c, recorded := newTestCLI(t, http.StatusAccepted, "")

		if err := c.deleteRepository("PRJ", "my-repo"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertRequest(t, recorded, http.MethodDelete, "/api/1.0/projects/PRJ/repos/my-repo")
		if recorded.atlassianTok != "no-check" {
			t.Errorf("X-Atlassian-Token = %q", recorded.atlassianTok)
		}
	})

	t.Run("escapes path segments", func(t *testing.T) {
		c, recorded := newTestCLI(t, http.StatusAccepted, "")

		if err := c.deleteRepository("PR J", "my/repo"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if recorded.path != "/api/1.0/projects/PR J/repos/my/repo" {
			t.Errorf("decoded path = %q", recorded.path)
		}
	})

	t.Run("does not exist", func(t *testing.T) {
		c, _ := newTestCLI(t, http.StatusNoContent, "")

		err := c.deleteRepository("PRJ", "missing")
		if err == nil || !strings.Contains(err.Error(), "does not exist") {
			t.Fatalf("expected a does-not-exist error, got %v", err)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		c, _ := newTestCLI(t, http.StatusUnauthorized, errorJson)

		err := c.deleteRepository("PRJ", "my-repo")
		if err == nil {
			t.Fatal("expected an error")
		}
		if !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "not permitted") {
			t.Errorf("error does not carry status and body: %v", err)
		}
	})
}

func TestConfirmDeletion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "matching slug", input: "my-repo\n", wantErr: false},
		{name: "matching slug with whitespace", input: "  my-repo  \n", wantErr: false},
		{name: "matching slug without newline", input: "my-repo", wantErr: false},
		{name: "wrong slug", input: "other\n", wantErr: true},
		{name: "empty line", input: "\n", wantErr: true},
		{name: "end of input", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var prompt strings.Builder
			err := confirmDeletion(strings.NewReader(tt.input), &prompt, "PRJ", "my-repo", "My Repo")
			if tt.wantErr && err == nil {
				t.Fatal("expected an error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(prompt.String(), "PRJ/my-repo") || !strings.Contains(prompt.String(), "My Repo") {
				t.Errorf("prompt does not name the repository: %q", prompt.String())
			}
		})
	}
}

func assertRequest(t *testing.T, recorded *recordedRequest, method string, path string) {
	t.Helper()
	if recorded.method != method {
		t.Errorf("method = %q, want %q", recorded.method, method)
	}
	if recorded.path != path {
		t.Errorf("path = %q, want %q", recorded.path, path)
	}
	if !strings.HasPrefix(recorded.authorization, "Basic ") {
		t.Errorf("Authorization = %q, want basic auth", recorded.authorization)
	}
}

func assertJsonContentType(t *testing.T, recorded *recordedRequest) {
	t.Helper()
	if !strings.HasPrefix(recorded.contentType, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", recorded.contentType)
	}
}

// assertJsonEqual compares two values through their JSON encoding, which
// sidesteps the difference between Go types and decoded JSON types.
func assertJsonEqual(t *testing.T, want interface{}, got interface{}) {
	t.Helper()
	wantJson, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("unable to marshal expected value: %v", err)
	}
	gotJson, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("unable to marshal actual value: %v", err)
	}
	if string(wantJson) != string(gotJson) {
		t.Errorf("got %s, want %s", gotJson, wantJson)
	}
}
