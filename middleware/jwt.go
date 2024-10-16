package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type UserClaims struct {
	ID        uint64
	Email     string
	CreatedAt time.Time
	Type  string
}

type contextKey string

const userContextKey contextKey = "userClaims"

func AuthJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the "Authorization" header
		tokenString := r.Header.Get("AuthToken")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Validate the token
		token, err := ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// If valid, set the user ID or email to context (optional)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userClaims := UserClaims{
				ID:        uint64(claims["id"].(float64)), // Convert to uint64 if necessary
				Email:     claims["user"].(string),
				CreatedAt: time.Unix(int64(claims["created_at"].(float64)), 0), // Convert to time.Time
				Type: claims["type"].(string),
			}

			// Create a new context with the user info
			ctx := context.WithValue(r.Context(), userContextKey, userClaims)
			r = r.WithContext(ctx)
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func CreateJWT(email string, createdAt time.Time, id uint64,userType string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	key := os.Getenv("secretKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         id,
		"type":       userType,
		"user":       email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"created_at": createdAt,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println("Error Creating the token", err)
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	key := os.Getenv("secretKey")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
}
