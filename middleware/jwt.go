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
	ID        uint
	Email     string
	CreatedAt time.Time
	Type  string
}

type contextKey string

const userContextKey contextKey = "userClaims"

func AuthJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		tokenString := r.Header.Get("AuthToken")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		
		token, err := ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userClaims UserClaims

			
			if id, ok := claims["id"].(float64); ok {
				userClaims.ID = uint(id)
			} else {
				http.Error(w, "Invalid token payload: id", http.StatusUnauthorized)
				return
			}

			
			if email, ok := claims["user"].(string); ok {
				userClaims.Email = email
			} else {
				http.Error(w, "Invalid token payload: user", http.StatusUnauthorized)
				return
			}

			
			if createdAt, ok := claims["created_at"].(float64); ok {
				userClaims.CreatedAt = time.Unix(int64(createdAt), 0)
			} else {
				http.Error(w, "Invalid token payload: created_at", http.StatusUnauthorized)
				return
			}

			
			if userType, ok := claims["type"].(string); ok {
				userClaims.Type = userType
			} else {
				http.Error(w, "Invalid token payload: type", http.StatusUnauthorized)
				return
			}

			
			ctx := context.WithValue(r.Context(), userContextKey, userClaims)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		
		next.ServeHTTP(w, r)
	})
}

func CreateJWT(email string, createdAt time.Time, id uint,userType string) (string, error) {
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
		"created_at": createdAt.Unix(),
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
