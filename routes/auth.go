package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
)

func RegisterAuthRoutes(router *mux.Router,userHandler *handlers.UserHandler) {
	router.HandleFunc("/signup",userHandler.SignUp).Methods("POST")
	router.HandleFunc("/login",userHandler.Login).Methods("POST")
}