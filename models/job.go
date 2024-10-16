package models

import "time"

type Job struct {
	ID                uint   // Primary key
	Title             string `json:"title"`
	Description       string `json:"description"`
	PostedOn          time.Time
	TotalApplications int
	CompanyName       string `json:"companyName"`
	PostedByID        uint   // Foreign key referring to User ID
}

