package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			unauthorized(w, "missing bearer token")
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, getSecret,
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)
		if err != nil {
			unauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSecret(token *jwt.Token) (any, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func User(ctx context.Context) string {
	user, _ := ctx.Value(userKey).(string)
	return user
}

func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
