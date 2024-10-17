package middleware

import "net/http"

func AuthApplicant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(userContextKey).(UserClaims)
		if !ok  {
			http.Error(w, "Unauthorized: Could not retrieve user info applicant", http.StatusUnauthorized)
			return
		}
		// Extract necessary fields from userInfo
		// email, okEmail := userInfo["user"].(string)
		// id, okID := userInfo["id"].(uint)
		// createdAt, okCreatedAt := userInfo["created_at"].(time.Time)
		usertType:= userInfo.Type
		
		if usertType != "applicant" {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
