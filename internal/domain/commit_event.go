package domain

import "time"

// CommitEvent представляет событие коммита из GitLab webhook
type CommitEvent struct {
	RepositoryID   string
	RepositoryName string
	Branch         string
	Author         string
	AuthorEmail    string
	CommitHash     string
	CommitMsg      string
	Timestamp      time.Time
	CommitURL      string
}
