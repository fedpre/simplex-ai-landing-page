package database

import (
	"log"

	"github.com/federicopregnolato/simplexai-landing-page/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Submission{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func CreateSubmission(submission *models.Submission) error {
	return db.Create(submission).Error
}

func GetSubmissions() ([]models.Submission, error) {
	var submissions []models.Submission
	err := db.Order("created_at DESC").Find(&submissions).Error
	return submissions, err
}
