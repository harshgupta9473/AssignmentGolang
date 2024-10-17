package handlers

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/harshgupta9473/recruitmentManagement/middleware"
	"github.com/harshgupta9473/recruitmentManagement/models"
	"github.com/harshgupta9473/recruitmentManagement/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var usereq models.User
	err := json.NewDecoder(r.Body).Decode(&usereq)
	if err != nil {
		http.Error(w, "error occured 1 ", http.StatusInternalServerError)
		return
	}
	tmpuser, err := utils.FindUserByEmail(usereq.Email)
	if err != nil {
		http.Error(w, "error occured 2 ", http.StatusInternalServerError)
		return
	}
	if tmpuser == nil {
		
		err = utils.InsertIntoUser(&usereq)
		if err != nil {
			http.Error(w, "error occured 3 ", http.StatusInternalServerError)
			return
		}
	} else {
		utils.WriteJSON(w, http.StatusContinue, "you are already signup Login to access account")
		return
	}
	utils.WriteJSON(w,http.StatusOK,"successful signup login to access your account")

}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := utils.FindUserByEmail(loginReq.Email)
	if err != nil {
		http.Error(w, "Error occurred while fetching user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	log.Println(user.Password)
	log.Println(loginReq.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.CreateJWT(user.Email, user.CreatedAt, (user.ID),user.UserType)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (u *UserHandler)GetAllJobs(w http.ResponseWriter,r *http.Request){
	jobs,err:=utils.GetAllJobs()
	if err!=nil{
		http.Error(w,"error",http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w,http.StatusOK,jobs)
}
