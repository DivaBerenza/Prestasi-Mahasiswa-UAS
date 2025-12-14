package middleware

import (
	"context"
	"net/http"
	"strings"

	"UAS/app/utils"
)

type ctxKey string

const (
	UserIDKey ctxKey = "user_id"
	RoleIDKey ctxKey = "role_id"
	PermKey   ctxKey = "permissions"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)
		ctx = context.WithValue(ctx, PermKey, claims.Permissions)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
