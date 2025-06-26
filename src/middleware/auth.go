package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "your-secret-key"))
var adminUser = getEnv("ADMIN_USER", "admin")
var adminPassword = getEnv("ADMIN_PASSWORD", "password")

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Authenticate checks if the provided username and password are valid
func Authenticate(username, password string) bool {
	return username == adminUser && password == adminPassword
}

// GenerateToken generates a JWT token for the authenticated user
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses a JWT token and returns the token and claims
func ParseToken(tokenStr string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Parse the JWT token
		tokenStr := cookie.Value
		claims := &Claims{}
		token, err := ParseToken(tokenStr, claims)

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
