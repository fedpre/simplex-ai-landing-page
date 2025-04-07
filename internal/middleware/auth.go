package middleware

import (
	"net/http"

	"github.com/federicopregnolato/simplexai-landing-page/internal/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsAuthenticated(r) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
