package models

import "time"

type User struct {
    ID              uint      // Primary key
    Name            string    `json:"name"`
    Email           string    `json:"email"`
    Address         string    `json:"address"`
    UserType        string    `json:"userType"` // "Admin" or "Applicant"
    Password    string    `json:"-"`
    ProfileHeadline string    `json:"profileHeadline"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type Profile struct {
    ID                uint   // Primary key
    UserID            uint   // Foreign key referring to User ID
    ResumeFileAddress string `json:"resumeFileAddress"`
    Skills            string `json:"skills"`
    Education         string `json:"education"`
    Experience        string `json:"experience"`
    Name              string `json:"name"`
    Email             string `json:"email"`
    Phone             string `json:"phone"`
}