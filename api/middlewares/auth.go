package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/0x-ximon/portman/api/services"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			next.ServeHTTP(w, r)
			return
		}

		token := parts[1]
		ctx := context.WithValue(r.Context(), services.TokenKey{}, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
