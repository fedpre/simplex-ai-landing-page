package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/federicopregnolato/simplexai-landing-page/internal/config"
	"github.com/federicopregnolato/simplexai-landing-page/internal/database"
	"github.com/federicopregnolato/simplexai-landing-page/internal/handlers"
)

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filename))
}

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Get the project root directory
	projectRoot := getProjectRoot()
	log.Printf("Project root: %s", projectRoot)

	// Get the absolute paths to the templates and static directories
	templatesDir := filepath.Join(projectRoot, "templates")
	staticDir := filepath.Join(projectRoot, "static")

	// Load templates with absolute paths
	tmplMain, err := template.ParseFiles(filepath.Join(templatesDir, "main.html"))
	if err != nil {
		log.Fatalf("Error loading main template: %v", err)
	}

	tmplAdmin, err := template.ParseFiles(filepath.Join(templatesDir, "admin.html"))
	if err != nil {
		log.Fatalf("Error loading admin template: %v", err)
	}

	// Initialize handlers
	mainHandler := handlers.NewMainHandler(cfg, tmplMain)
	adminHandler := handlers.NewAdminHandler(cfg, tmplAdmin)

	// Create router
	mux := http.NewServeMux()

	// Register routes
	mainHandler.RegisterRoutes(mux)
	adminHandler.RegisterRoutes(mux)

	// Serve static files
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
