package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/sensetion/tgGitlabBot/internal/controller/http/response"
)

func WebhookAuth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Gitlab-Token")

			if token == "" {
				response.Error(w, http.StatusUnauthorized, "missing webhook token")
				return
			}

			// Безопасное сравнение строк (защита от timing attacks)
			if subtle.ConstantTimeCompare([]byte(token), []byte(secret)) != 1 {
				response.Error(w, http.StatusUnauthorized, "invalid webhook token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
