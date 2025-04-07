package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/federicopregnolato/simplexai-landing-page/internal/config"
	"github.com/federicopregnolato/simplexai-landing-page/internal/database"
	"github.com/federicopregnolato/simplexai-landing-page/internal/middleware"
	"github.com/federicopregnolato/simplexai-landing-page/internal/models"
	"github.com/federicopregnolato/simplexai-landing-page/internal/utils"
)

type AdminHandler struct {
	config    *config.Config
	tmplAdmin *template.Template
}

func NewAdminHandler(cfg *config.Config, tmplAdmin *template.Template) *AdminHandler {
	return &AdminHandler{
		config:    cfg,
		tmplAdmin: tmplAdmin,
	}
}

func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /admin", h.handleAdminPage)
	mux.HandleFunc("POST /admin/login", middleware.CSRFMiddleware(h.handleLogin))
	mux.HandleFunc("POST /admin/logout", middleware.CSRFMiddleware(h.handleLogout))
}

func (h *AdminHandler) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	isAuth := utils.IsAuthenticated(r)

	var csrfToken string
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil || csrfCookie.Value == "" {
		csrfToken = utils.GenerateCSRFToken()
		http.SetCookie(w, utils.NewSecureCookie("csrf_token", csrfToken, time.Now().Add(h.config.CookieDuration)))
	} else {
		csrfToken = csrfCookie.Value
	}

	// Get error message from query parameter
	loginError := ""
	if r.URL.Query().Get("error") == "1" {
		loginError = "Invalid email or password"
	}

	if !isAuth {
		data := models.AdminPageData{
			IsAuthenticated: false,
			CSRFToken:       csrfToken,
			LoginError:      loginError,
		}
		h.tmplAdmin.Execute(w, data)
		return
	}

	// Get submissions from database
	submissions, err := database.GetSubmissions()
	if err != nil {
		http.Error(w, "Error fetching submissions", http.StatusInternalServerError)
		return
	}

	data := models.AdminPageData{
		IsAuthenticated: true,
		CSRFToken:       csrfToken,
		Submissions:     submissions,
	}
	h.tmplAdmin.Execute(w, data)
}

func (h *AdminHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == h.config.AdminEmail && password == h.config.AdminPassword {
		http.SetCookie(w, utils.NewSecureCookie("auth", "authenticated", time.Now().Add(h.config.CookieDuration)))
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin?error=1", http.StatusSeeOther)
}

func (h *AdminHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
