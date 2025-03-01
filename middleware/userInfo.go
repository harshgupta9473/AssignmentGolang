package middleware

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/harshgupta9473/recruitmentManagement/config"
	"github.com/harshgupta9473/recruitmentManagement/models"
)

func InfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
		userInfo, ok := r.Context().Value(userContextKey).(UserClaims)
		if !ok  {
			http.Error(w, "Unauthorized: Could not retrieve user info 2", http.StatusUnauthorized)
			return
		}
		
		

		if !FindUser(userInfo.Email, userInfo.ID, userInfo.Type) {
			http.Error(w, "Unauthorized 1", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func FindUser(email string, id uint,  userType string) bool {
	db := config.GetDB()
	query := `SELECT id, name,email FROM users WHERE id = $1 AND email = $2 AND user_type=$3`
	result := db.QueryRow(query, id, email,userType)

	var user models.User 
	err := result.Scan(&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			
			return false
		}
		
		log.Println("Error querying user:", err)
		return false
	}
	return true 
}


func ExtractUserClaimsFromContext(r *http.Request)(*UserClaims,error) {
	userInfo, ok := r.Context().Value(userContextKey).(UserClaims)
	if !ok {
		
		return nil,fmt.Errorf("Unauthorized: Could not retrieve user info applicant")
	}
	return &userInfo,nil
}
