package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
)

func RegisterAuthRoutes(router *mux.Router, userHandler *handlers.UserHandler) {
	router.HandleFunc("/signup", userHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
}
