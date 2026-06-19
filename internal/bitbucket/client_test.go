package bitbucket

import (
	"context"
	"errors"
	"testing"
)

type MockClient struct {
	MockRepo  Repository
	MockError error
}

func (m *MockClient) GetRepository(ctx context.Context, workspace string, repoSlug string) (Repository, error) {
	return m.MockRepo, m.MockError
}

func TestGetRepository(t *testing.T) {
	tests := []struct {
		name          string
		mockClient    *MockClient
		workspace     string
		repoSlug      string
		expectedError error
		expectedName  string
	}{
		{
			name: "Successful Repository Fetch",
			mockClient: &MockClient{
				MockRepo:  Repository{Name: "mcp-demo", Description: "A test repo", IsPrivate: true},
				MockError: nil,
			},
			workspace:     "my-workspace",
			repoSlug:      "mcp-demo",
			expectedError: nil,
			expectedName:  "mcp-demo",
		},
		{
			name: "Repository Not Found",
			mockClient: &MockClient{
				MockRepo:  Repository{},
				MockError: errors.New("not found"),
			},
			workspace:     "my-workspace",
			repoSlug:      "unknown-repo",
			expectedError: errors.New("not found"),
			expectedName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo, err := tt.mockClient.GetRepository(ctx, tt.workspace, tt.repoSlug)

			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if repo.Name != tt.expectedName {
				t.Errorf("expected repo name %s, got %s", tt.expectedName, repo.Name)
			}
		})
	}
}
