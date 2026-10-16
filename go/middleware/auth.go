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
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || raw == "" {
			unauthorized(w, "missing bearer token")
			return
		}

		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(raw, claims,
			func(*jwt.Token) (any, error) { return []byte(os.Getenv("JWT_SECRET")), nil },
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
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

func User(ctx context.Context) string {
	user, _ := ctx.Value(userKey).(string)
	return user
}
func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
