package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"warehouse-api/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Generate JWT token
func GenerateToken(userID int, username, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key"
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Verify JWT token
func VerifyToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("AuthMiddleware - Path: %s, Method: %s\n", r.URL.Path, r.Method)

		authHeader := r.Header.Get("Authorization")
		fmt.Printf("AuthMiddleware - Authorization Header: %s\n", authHeader)

		if authHeader == "" {
			fmt.Println("AuthMiddleware - No Authorization header")
			sendJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Authorization header required",
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Printf("AuthMiddleware - Invalid format. Parts: %d, First: %s\n", len(parts), parts[0])
			sendJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
			return
		}

		token := parts[1]
		fmt.Printf("AuthMiddleware - Token extracted: %s...\n", token[:20])

		claims, err := VerifyToken(token)
		if err != nil {
			fmt.Printf("AuthMiddleware - Token verification failed: %v\n", err)
			sendJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			})
			return
		}

		fmt.Printf("AuthMiddleware - Token verified. User: %s, Role: %s\n", claims.Username, claims.Role)

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Role-based authorization middleware
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*Claims)
			if !ok {
				sendJSON(w, http.StatusUnauthorized, models.ErrorResponse{
					Success: false,
					Message: "Unauthorized",
				})
				return
			}

			// Check if user's role is in allowed roles
			allowed := false
			for _, role := range roles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				sendJSON(w, http.StatusForbidden, models.ErrorResponse{
					Success: false,
					Message: "Insufficient permissions",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Get user from context
func GetUserFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return claims, nil
}

// CORS middleware
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper function to send JSON response
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
