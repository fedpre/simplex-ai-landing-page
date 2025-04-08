package models

import "time"

type MainPageData struct {
	AppTitle               string
	SubmissionConfirmation bool
}

type AdminPageData struct {
	IsAuthenticated bool
	Submissions     []Submission
	LoginError      string
}

type Submission struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
