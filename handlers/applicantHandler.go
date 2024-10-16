package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/harshgupta9473/recruitmentManagement/models"
	"github.com/harshgupta9473/recruitmentManagement/utils"
	"github.com/joho/godotenv"
)

type ApplicantHandler struct {
	DB *sql.DB
}

func (applicant *ApplicantHandler) HandleUpload(w http.ResponseWriter, r http.Request) {
	err := r.ParseMultipartForm(10 << 22)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(handler.Filename, "/pdf") && !strings.HasSuffix(handler.Filename, ".docx") {
		http.Error(w, "Invalid file type. Only PDF and DOCX are allowed.", http.StatusBadRequest)
		return
	}

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	err = godotenv.Load()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	resumeAPIURL := os.Getenv("resumeAPIURL")
	apiKey := os.Getenv("apiKey")

	req, err := http.NewRequest("POST", resumeAPIURL, &fileBuffer)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request to resume parser API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var userProfile models.Profile
	err = json.NewDecoder(resp.Body).Decode(&userProfile)
	if err != nil {
		http.Error(w, "Error parsing response from resume parser API", http.StatusInternalServerError)
		return
	}
}

func (applicant *ApplicantHandler) ApplyForJobByJobID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jobid := query.Get("job_id") // Retrieves the value of job_id

	if jobid == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}
	jobID, err := strconv.ParseUint(jobid, 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	idstring:=r.Header.Get("id")
	id, err := strconv.ParseUint(idstring, 10, 32)
	if err != nil {
		http.Error(w, "Invalid applicant ID", http.StatusBadRequest)
		return
	}
	err=utils.ApplyForJob(uint(jobID),uint(id))
	if err!=nil{
		http.Error(w,"not able to apply for the job",http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w,http.StatusOK,"job application submitted successfully")
}
