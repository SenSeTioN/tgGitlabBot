package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sensetion/tgGitlabBot/internal/adapters/http/handler"
	"github.com/sensetion/tgGitlabBot/internal/config"
)

func Init(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	setupRouter(r, cfg)
	setupHandlers(r)

	return r
}

func setupHandlers(r *chi.Mux) {
	healthHandler := handler.NewHealthHandler(nil)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	// 404 handler
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
	r.Use(middleware.Compress(5))
}
