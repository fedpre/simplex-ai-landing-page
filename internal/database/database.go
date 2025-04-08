package database

import (
	"fmt"
	"log"
	"os"

	"github.com/federicopregnolato/simplexai-landing-page/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *gorm.DB

func InitDB() {
	var err error
	dbUrl := os.Getenv("TURSO_URL")
	dbToken := os.Getenv("TURSO_TOKEN")
	url := fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken)

	sqliteCfg := sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}

	// db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	db, err = gorm.Open(sqlite.New(sqliteCfg), &gorm.Config{})

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
