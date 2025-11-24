package gitlab

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sensetion/tgGitlabBot/internal/domain"
	"github.com/sensetion/tgGitlabBot/pkg/logger"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

type pushEventPayload struct {
	ObjectKind   string      `json:"object_kind"`
	Ref          string      `json:"ref"`
	UserUsername string      `json:"user_username"`
	UserEmail    string      `json:"user_email"`
	Project      projectInfo `json:"project"`
	Commits      []commit    `json:"commits"`
}

type projectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
}

type commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url"`
	Author    author    `json:"author"`
}

type author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *Parser) ParsePushEvent(payload []byte) (*domain.CommitEvent, error) {
	var event pushEventPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.PrettyStructurePrint("Event BODY :", event)

	if event.ObjectKind != "push" {
		return nil, fmt.Errorf("unsupported object_kind: %s", event.ObjectKind)
	}

	if len(event.Commits) == 0 {
		return nil, fmt.Errorf("no commits found in payload")
	}

	lastCommit := event.Commits[len(event.Commits)-1]

	branch := p.extractBranch(event.Ref)

	return &domain.CommitEvent{
		RepositoryID:   fmt.Sprintf("%d", event.Project.ID),
		RepositoryName: event.Project.PathWithNamespace,
		Branch:         branch,
		Author:         lastCommit.Author.Name,
		AuthorEmail:    lastCommit.Author.Email,
		CommitHash:     lastCommit.ID,
		CommitMsg:      lastCommit.Message,
		Timestamp:      lastCommit.Timestamp,
		WebURL:         event.Project.WebURL,
		CommitURL:      lastCommit.URL,
	}, nil
}

func (p *Parser) extractBranch(ref string) string {
	if len(ref) > 11 && ref[:11] == "refs/heads/" {
		return ref[11:]
	}
	return ref
}
