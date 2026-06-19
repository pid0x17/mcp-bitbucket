package bitbucket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClient_GetRepository(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/repositories/my-workspace/mcp-demo"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if auth := r.Header.Get("Authorization"); auth != "Bearer fake_test_token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("expected Accept header to be application/json, got %s", accept)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "mcp-demo",
			"description": "A test repository",
			"is_private": true
		}`))
	}))

	defer mockServer.Close()

	client := NewHTTPClient(
		WithBaseURL(mockServer.URL),
		WithToken("fake_test_token"),
	)

	ctx := context.Background()
	repo, err := client.GetRepository(ctx, "my-workspace", "mcp-demo")

	if err != nil {
		t.Fatalf("did not expect an error, got: %v", err)
	}

	if repo.Name != "mcp-demo" {
		t.Errorf("expected repo name 'mcp-demo', got '%s'", repo.Name)
	}
	if !repo.IsPrivate {
		t.Errorf("expected repo to be private")
	}
}
