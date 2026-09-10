package cli_test

import (
	cli "github.com/swisscom/bitbucket-cli/internal"
	"testing"
)

func TestParseBitbucketRemote(t *testing.T) {
	tests := []struct {
		name           string
		remoteURL      string
		wantProjectKey string
		wantSlug       string
		wantBaseURL    string
		wantErr        bool
	}{
		// ── HTTPS, root deployment ────────────────────────────────────────────
		{
			name:           "https root with .git",
			remoteURL:      "https://bitbucket.example.com/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com/rest",
		},
		{
			name:           "https root without .git",
			remoteURL:      "https://bitbucket.example.com/scm/MYPROJ/my-repo",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com/rest",
		},
		{
			name:           "https root with embedded username",
			remoteURL:      "https://jsmith@bitbucket.example.com/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com/rest",
		},
		{
			name:           "https root with non-standard port",
			remoteURL:      "https://bitbucket.example.com:8443/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com:8443/rest",
		},

		// ── HTTPS, non-root (context path) deployment ────────────────────────
		{
			name:           "https with single context path segment",
			remoteURL:      "https://bitbucket.example.com/bitbucket/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com/bitbucket/rest",
		},
		{
			name:           "https with nested context path",
			remoteURL:      "https://example.com/tools/bitbucket/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://example.com/tools/bitbucket/rest",
		},
		{
			name:           "https context path with username and port",
			remoteURL:      "https://jsmith@bitbucket.example.com:8443/stash/scm/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "https://bitbucket.example.com:8443/stash/rest",
		},

		// ── SSH URL format ────────────────────────────────────────────────────
		{
			name:           "ssh with port",
			remoteURL:      "ssh://git@bitbucket.example.com:7999/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},
		{
			name:           "ssh without port",
			remoteURL:      "ssh://git@bitbucket.example.com/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},
		{
			name:           "ssh without .git suffix",
			remoteURL:      "ssh://git@bitbucket.example.com:7999/MYPROJ/my-repo",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},

		// ── git+ssh URL format ────────────────────────────────────────────────
		{
			name:           "git+ssh with port",
			remoteURL:      "git+ssh://git@bitbucket.example.com:7999/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},
		{
			name:           "git+ssh without port",
			remoteURL:      "git+ssh://git@bitbucket.example.com/MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},

		// ── SCP-style ─────────────────────────────────────────────────────────
		{
			name:           "scp with user",
			remoteURL:      "git@bitbucket.example.com:MYPROJ/my-repo.git",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},
		{
			name:           "scp without .git suffix",
			remoteURL:      "git@bitbucket.example.com:MYPROJ/my-repo",
			wantProjectKey: "MYPROJ",
			wantSlug:       "my-repo",
			wantBaseURL:    "",
		},

		// ── Error cases ───────────────────────────────────────────────────────
		{
			name:      "https missing /scm/ segment",
			remoteURL: "https://bitbucket.example.com/MYPROJ/my-repo.git",
			wantErr:   true,
		},
		{
			name:      "ssh missing repo segment",
			remoteURL: "ssh://git@bitbucket.example.com:7999/MYPROJ",
			wantErr:   true,
		},
		{
			name:      "scp missing repo segment",
			remoteURL: "git@bitbucket.example.com:MYPROJ",
			wantErr:   true,
		},
		{
			name:      "unrecognised URL with no colon",
			remoteURL: "not-a-url",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectKey, slug, baseURL, err := cli.ParseBitbucketRemote(tt.remoteURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got projectKey=%q slug=%q baseURL=%q", projectKey, slug, baseURL)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if projectKey != tt.wantProjectKey {
				t.Errorf("projectKey = %q, want %q", projectKey, tt.wantProjectKey)
			}
			if slug != tt.wantSlug {
				t.Errorf("slug = %q, want %q", slug, tt.wantSlug)
			}
			if baseURL != tt.wantBaseURL {
				t.Errorf("baseURL = %q, want %q", baseURL, tt.wantBaseURL)
			}
		})
	}
}
