package models

type JobLisitingReq struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	CompanyName       string `json:"companyName"`
	PostedByID        uint  
}

type JobOpeningResponse struct {
	Job
	Applicants []User
}

type ResumeResponse struct {
	Education []struct {
		Name string      `json:"name"`
		URL  string       `json:"url"`
	}                     `json:"education"`
	Email      string     `json:"email"`
	Experience []struct {
		Dates []string      `json:"dates"`
		Name  string        `json:"name"`
		URL   string         `json:"url"`
	}                     `json:"experience"`
	Name   string   `json:"name"`
	Phone  string   `json:"phone"`
	Skills []string `json:"skills"`
}
