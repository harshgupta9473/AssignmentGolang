package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
	"github.com/harshgupta9473/recruitmentManagement/middleware"
)

func RegisterCommonRoutes(router *mux.Router,uh *handlers.UserHandler ) {
	subrouter:=router.NewRoute().Subrouter()
	subrouter.Use(middleware.AuthJwt)
	subrouter.Use(middleware.InfoMiddleware)

	subrouter.HandleFunc("/jobs",uh.GetAllJobs).Methods(http.MethodGet)
}
