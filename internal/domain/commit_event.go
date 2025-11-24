package domain

import "time"

type CommitEvent struct {
	RepositoryID   string
	RepositoryName string
	Branch         string
	Author         string
	AuthorEmail    string
	CommitHash     string
	CommitMsg      string
	Timestamp      time.Time
	WebURL         string
	CommitURL      string
}
