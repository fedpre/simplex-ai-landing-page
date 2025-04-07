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

type MainHandler struct {
	config   *config.Config
	tmplMain *template.Template
}

func NewMainHandler(cfg *config.Config, tmplMain *template.Template) *MainHandler {
	return &MainHandler{
		config:   cfg,
		tmplMain: tmplMain,
	}
}

func (h *MainHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.handleMainPage)
	mux.HandleFunc("POST /submit", middleware.CSRFMiddleware(h.handleSubmit))
}

func (h *MainHandler) handleMainPage(w http.ResponseWriter, r *http.Request) {
	var csrfToken string
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil || csrfCookie.Value == "" {
		csrfToken = utils.GenerateCSRFToken()
		http.SetCookie(w, utils.NewSecureCookie("csrf_token", csrfToken, time.Now().Add(h.config.CookieDuration)))
	} else {
		csrfToken = csrfCookie.Value
	}

	data := models.MainPageData{
		AppTitle:               "SimplexAI",
		SubmissionConfirmation: false,
		CSRFToken:              csrfToken,
	}
	h.tmplMain.Execute(w, data)
}

func (h *MainHandler) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	// Create submission in database
	submission := &models.Submission{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.CreateSubmission(submission); err != nil {
		http.Error(w, "Error saving submission", http.StatusInternalServerError)
		return
	}

	// Generate new CSRF token for the next form submission
	csrfToken := utils.GenerateCSRFToken()
	http.SetCookie(w, utils.NewSecureCookie("csrf_token", csrfToken, time.Now().Add(h.config.CookieDuration)))

	// Show success message on the main page
	data := models.MainPageData{
		AppTitle:               "SimplexAI",
		SubmissionConfirmation: true,
		CSRFToken:              csrfToken,
	}
	h.tmplMain.Execute(w, data)
}
