package models

import "time"

type User struct {
    ID              uint      `json:"id"`
    Name            string    `json:"name"`
    Email           string    `json:"email"`
    Address         string    `json:"address"`
    UserType        string    `json:"userType"` 
    Password    string        `json:"password"`
    ProfileHeadline string    `json:"profileHeadline"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}

type Profile struct {
    ID                uint    
    UserID            uint   
    ResumeFileAddress string `json:"resumeFileAddress"`
    Skills            string `json:"skills"`
    Education         string `json:"education"`
    Experience        string `json:"experience"`
    Name              string `json:"name"`
    Email             string `json:"email"`
    Phone             string `json:"phone"`
}