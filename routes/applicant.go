package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
	"github.com/harshgupta9473/recruitmentManagement/middleware"
)

func RegisterApplicantRoutes(router *mux.Router, ap *handlers.ApplicantHandler) {
	subrouter:=router.NewRoute().Subrouter()
	subrouter.Use(middleware.AuthJwt)
	subrouter.Use(middleware.InfoMiddleware)
	subrouter.Use(middleware.AuthApplicant)

	subrouter.HandleFunc("/jobs/apply",ap.ApplyForJobByJobID).Methods(http.MethodGet)
	subrouter.HandleFunc("/uploadResume",ap.HandleUpload).Methods(http.MethodPost)
}
