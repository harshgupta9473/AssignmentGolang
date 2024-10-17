package handlers

import (
	"bytes"
	"database/sql"
	"log"
	"time"

	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"github.com/harshgupta9473/recruitmentManagement/middleware"
	"github.com/harshgupta9473/recruitmentManagement/models"
	"github.com/harshgupta9473/recruitmentManagement/utils"
	"github.com/joho/godotenv"
)

type ApplicantHandler struct {
	DB *sql.DB
}

func NewApplicantHandler(db *sql.DB) *ApplicantHandler {
	return &ApplicantHandler{
		DB: db,
	}
}

func (applicant *ApplicantHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 22)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	userInfo, err := middleware.ExtractUserClaimsFromContext(r)
	if err != nil {
		http.Error(w, "Unable to authorise user", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(handler.Filename, ".pdf") && !strings.HasSuffix(handler.Filename, ".docx") {
		http.Error(w, "Invalid file type. Only PDF and DOCX are allowed.", http.StatusBadRequest)
		return
	}

	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	err = godotenv.Load()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	uploadedURL, err := utils.SaveFile(handler)
	if err!=nil{
		http.Error(w,"error uploading resume",http.StatusInternalServerError)
		return
	}
	resumeAPIURL := os.Getenv("resumeAPIURL")
	apiKey := os.Getenv("apiKey")

	req, err := http.NewRequest("POST", resumeAPIURL, fileBuffer)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second, 
	}
	resumeResp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request to resume parser API", http.StatusInternalServerError)
		return
	}
	defer resumeResp.Body.Close()
	var userProfile models.ResumeResponse
	err = json.NewDecoder(resumeResp.Body).Decode(&userProfile)
	if err != nil {
		http.Error(w, "Error parsing response from resume parser API", http.StatusInternalServerError)
		return
	}

	log.Println(userProfile)

	education := utils.ExtractFromField(userProfile.Education, "name")
	experience := utils.ExtractFromField(userProfile.Experience, "name")
	skills := strings.Join(userProfile.Skills, ", ")
	err = utils.InsertIntoProfile(userInfo.ID, uploadedURL, skills, education, experience, userProfile.Name, userProfile.Email, userProfile.Phone)
	if err != nil {
		log.Println(err)
		http.Error(w, "error storing resume in database", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Resume Uploaded Successfully")
}

func (applicant *ApplicantHandler) ApplyForJobByJobID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jobid := query.Get("job_id") 
	if jobid == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}
	jobID, err := strconv.ParseUint(jobid, 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	userInfo, err := middleware.ExtractUserClaimsFromContext(r)
	if err != nil {
		http.Error(w, "error accessing data abbout user", http.StatusInternalServerError)
		return
	}
	err = utils.ApplyForJob(uint(jobID), userInfo.ID)
	if err != nil {
		http.Error(w, "not able to apply for the job", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "job application submitted successfully")
}
