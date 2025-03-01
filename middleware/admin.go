package middleware

import (
	"net/http"
)

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(userContextKey).(UserClaims)
		if !ok  {
			http.Error(w, "Unauthorized: Could not retrieve user info admin", http.StatusUnauthorized)
			return
		}
		// email, okEmail := userInfo["user"].(string)
		// id, okID := userInfo["id"].(uint)
		// createdAt, okCreatedAt := userInfo["created_at"].(time.Time)
		usertType:= userInfo.Type

		
		if usertType!="admin"||usertType==""{
			http.Error(w,"not allowed",http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w,r)
	})
}
