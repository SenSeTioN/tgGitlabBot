package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sensetion/tgGitlabBot/internal/controller/http/handler"
	chimw "github.com/sensetion/tgGitlabBot/internal/controller/http/middleware"
	"github.com/sensetion/tgGitlabBot/pkg/config"
)

func Init(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	setupRouter(r, cfg)
	setupHandlers(r, cfg)

	return r
}

func setupHandlers(r *chi.Mux, cfg *config.Config) {
	healthHandler := handler.NewHealthHandler(nil)
	webhookHandler := handler.NewWebhookHandler()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GitLab Telegram Bot API"))
	})

	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	r.Route("/webhook", func(wr chi.Router) {
		wr.Use(chimw.WebhookAuth(cfg.GitLab.WebhookSecret))
		wr.Use(middleware.AllowContentType("application/json"))

		wr.Post("/gitlab", webhookHandler.HandleGitLabPush)
	})

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"route not found"}`))
	})
}

func setupRouter(r *chi.Mux, cfg *config.Config) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Server.ReadTimeout))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Compress(cfg.Server.CompressSize))
}
