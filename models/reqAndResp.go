package models




type JobLisitingReq struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	TotalApplications int    `json:"totalApplicants"`
	CompanyName       string `json:"companyName"`
	PostedByID        uint   // Foreign key referring to User ID
}

type JobOpeningResponse struct{
	Job
	Applicants []User
}