package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/harshgupta9473/recruitmentManagement/config"
	"github.com/harshgupta9473/recruitmentManagement/models"
)

func InfoMiddleware(next http.Handler) http.Handler {
	const userContextKey = "userInfo"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user info from the context
		userInfo, ok := r.Context().Value(userContextKey).(map[string]interface{})
		if !ok || userInfo == nil {
			http.Error(w, "Unauthorized: Could not retrieve user info", http.StatusUnauthorized)
			return
		}
		// Extract necessary fields from userInfo
		email, okEmail := userInfo["user"].(string)
		id, okID := userInfo["id"].(uint)
		createdAt, okCreatedAt := userInfo["created_at"].(time.Time)
		usertType, oktype := userInfo["type"].(string)

		if !okEmail || !okID || !okCreatedAt || !oktype {
			http.Error(w, "Invalid user info", http.StatusUnauthorized)
			return
		}

		// Check if user exists
		if !FindUser(email, id, createdAt, usertType) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func FindUser(email string, id uint, createdAt time.Time, userType string) bool {
	db := config.GetDB()
	query := `SELECT id, name,email, FROM users WHERE id = $1 AND email = $2 AND user_type=$3 AND created_at = $4`
	result := db.QueryRow(query, id, email, createdAt, userType)

	var user models.User // Assuming you have a User struct defined
	err := result.Scan(&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return false
		}
		// Log any other error
		log.Println("Error querying user:", err)
		return false
	}
	return true // User exists
}
