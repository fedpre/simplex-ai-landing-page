package models

import "time"

type MainPageData struct {
	AppTitle               string
	SubmissionConfirmation bool
	CSRFToken              string
}

type AdminPageData struct {
	IsAuthenticated bool
	CSRFToken       string
	Submissions     []Submission
	LoginError      string
}

type Submission struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
