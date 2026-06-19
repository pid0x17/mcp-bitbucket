package bitbucket

import "context"

type RepositoryFetcher interface {
	GetRepository(ctx context.Context, workspace string, repoSlug string) (Repository, error)
}
