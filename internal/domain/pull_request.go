package domain

type PRStatus string

const (
	PRStatusOpen   PRStatus = "open"
	PRStatusMerged PRStatus = "merged"
)

type PullRequest struct {
	ID            string
	Title         string
	Description   string
	AuthorID      string
	TeamID        string
	Reviewers     []string
	Status        PRStatus
	MergedAt      *string
}
