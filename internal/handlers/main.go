package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/federicopregnolato/simplexai-landing-page/internal/config"
	"github.com/federicopregnolato/simplexai-landing-page/internal/database"
	"github.com/federicopregnolato/simplexai-landing-page/internal/models"
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
	mux.HandleFunc("POST /submit", h.handleSubmit)
}

func (h *MainHandler) handleMainPage(w http.ResponseWriter, r *http.Request) {
	data := models.MainPageData{
		AppTitle:               "SimplexAI",
		SubmissionConfirmation: false,
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

	// Create submission in database
	submission := &models.Submission{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.CreateSubmission(submission); err != nil {
		http.Error(w, "Error saving submission", http.StatusInternalServerError)
		return
	}

	// Show success message on the main page
	data := models.MainPageData{
		AppTitle:               "SimplexAI",
		SubmissionConfirmation: true,
	}
	h.tmplMain.Execute(w, data)
}
