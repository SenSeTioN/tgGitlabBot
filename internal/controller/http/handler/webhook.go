package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/sensetion/tgGitlabBot/internal/adapter/gitlab"
	"github.com/sensetion/tgGitlabBot/internal/controller/http/response"
	"github.com/sensetion/tgGitlabBot/pkg/logger"
)

type WebhookHandler struct {
	parser *gitlab.Parser
}

func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		parser: gitlab.NewParser(),
	}
}

func (h *WebhookHandler) HandleGitLabPush(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get("X-Gitlab-Event")
	if eventType != "Push Hook" {
		log.Printf("ℹ️ Ignoring event type: %s", eventType)
		response.JSON(w, http.StatusOK, map[string]string{"status": "ignored"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read body: %v", err)
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	defer r.Body.Close()

	event, err := h.parser.ParsePushEvent(body)
	if err != nil {
		log.Printf("❌ Parse error: %v", err)
		response.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	log.Println("🚀 GitLab Push Event Received:")
	logger.PrettyStructurePrint("Event :", event)

	response.JSON(w, http.StatusOK, map[string]string{"status": "processed"})
}
