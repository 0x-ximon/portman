package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
)

type middleware func(http.Handler) http.Handler

type Middleware struct {
	DbConn *pgx.Conn
}

func (m *Middleware) NewChain(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}

func (m *Middleware) ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo := repositories.New(m.DbConn)
		ctx := r.Context()

		apiKey := r.Header.Get("X-Api-Key")
		if apiKey != "" {
			user, err := repo.FindUserByApiKey(ctx, &apiKey)
			if err == nil {
				ctx := context.WithValue(ctx, services.ClaimsKey{}, &services.Claims{
					ID:            user.ID,
					EmailAddress:  user.EmailAddress,
					WalletAddress: user.WalletAddress,
				})

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

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

		claims, err := services.ValidateJWT(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, services.ClaimsKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
