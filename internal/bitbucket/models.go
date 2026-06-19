package bitbucket

/**
 * Repository represents a read-only view of a Bitbucket repository.
 */
type Repository struct {
	Name        string
	Description string
	IsPrivate   bool
}
