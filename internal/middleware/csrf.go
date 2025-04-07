package middleware

import (
	"net/http"

	"github.com/federicopregnolato/simplexai-landing-page/internal/utils"
)

func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if !utils.ValidateCSRFToken(r) {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
