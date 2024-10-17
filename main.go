package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/recruitmentManagement/database"
	"github.com/harshgupta9473/recruitmentManagement/handlers"
	"github.com/harshgupta9473/recruitmentManagement/routes"
	"github.com/harshgupta9473/recruitmentManagement/utils"
)

func main() {
	db:=database.InitDB()
	database.InitTable(db)
	utils.InitAWS()
	router:=mux.NewRouter()
	

	userHandler:=handlers.NewUserHandler(db)
	routes.RegisterAuthRoutes(router,userHandler)

	applicantHandler:=handlers.NewApplicantHandler(db)
	routes.RegisterApplicantRoutes(router,applicantHandler)

	adminHandler:=handlers.NewAdminHandler(db)
	routes.RegisterAdminRoutes(router,adminHandler)

	routes.RegisterCommonRoutes(router,userHandler)



	s:=&http.Server{
		Addr: ":4000",
		Handler: router,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func(){
		err:=s.ListenAndServe()
		if err!=nil{
			log.Fatal(err)
		}
	}()

	sigChan:=make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt)
	signal.Notify(sigChan,os.Kill)

	sig:=<-sigChan
	log.Println("recieved terminate, graceful shutdown",sig)

	tc,_:=context.WithTimeout(context.Background(),30*time.Second)

	s.Shutdown(tc)

}