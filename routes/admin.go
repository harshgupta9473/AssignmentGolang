package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
	"github.com/harshgupta9473/recruitmentManagement/middleware"
)

func RegisterAdminRoutes(router *mux.Router,ad *handlers.AdminHandler) {
	subrouter:=router.PathPrefix("/admin").Subrouter()
	subrouter.Use(middleware.AuthJwt)
	subrouter.Use(middleware.InfoMiddleware)
	subrouter.Use(middleware.AuthAdmin)

	subrouter.HandleFunc("/job/{job_id}",ad.FetchAboutJobOpeninBYID).Methods(http.MethodGet)
	subrouter.HandleFunc("/applicants",ad.FetchAllUsers).Methods(http.MethodGet)
	subrouter.HandleFunc("/applicants/{applicant_id}",ad.FetchaUserByID).Methods(http.MethodGet)
	

	subrouter.HandleFunc("/job",ad.CreateJobOpenings).Methods(http.MethodPost)
}