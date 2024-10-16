package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/models"
	"github.com/harshgupta9473/recruitmentManagement/utils"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{
		db: db,
	}
}

func (AH *AdminHandler) CreateJobOpenings(w http.ResponseWriter, r *http.Request) {
	var jobReq models.JobLisitingReq
	err := json.NewDecoder(r.Body).Decode(&jobReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = utils.InsertIntoJobs(&jobReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (AH *AdminHandler) FetchAboutJobOpeninBYID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["job_id"]

	jobIDInt, err := strconv.ParseUint(jobID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	job,err:=utils.FindJobByID(uint(jobIDInt))
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	if job==nil{
		http.Error(w,"job does not exists",http.StatusBadRequest)
		return
	}
	applicants,err:=utils.GetApplicantsForJob(job.ID)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	resp:=&models.JobOpeningResponse{
		Job: *job,
		Applicants: applicants,
	}
	utils.WriteJSON(w,http.StatusOK,resp)
}

func (AH *AdminHandler)FetchAllUsers(w http.ResponseWriter, r *http.Request){
	users,err:=utils.FetchAllUSers()
	if err!=nil{
		http.Error(w,"internal error",http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w,http.StatusOK,users)
}

func (AH *AdminHandler)FetchaUserByID(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	id:=vars["applicant_id"]

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	user,err:=utils.GetUSerBYID(uint(userID))
	if err!=nil{
		http.Error(w,"internal error",http.StatusInternalServerError)
		return
	}
	if user==nil{
		http.Error(w,"user not found",http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w,http.StatusOK,user)
}