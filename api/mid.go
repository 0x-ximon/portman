package main

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/0x-ximon/portman/api/services"
	"github.com/google/uuid"
)

type middleware func(http.Handler) http.Handler

type Middleware struct{}

func (m *Middleware) NewChain(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}

		return next
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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

func (m *Middleware) Logging(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(services.GetLogWriter(), nil))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		id := uuid.New().String()
		reqLogger := logger.With(
			slog.String("id", id),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		ctx := context.WithValue(r.Context(), services.LoggerKey{}, reqLogger)
		next.ServeHTTP(w, r.WithContext(ctx))

		reqLogger.Info("request completed",
			slog.Duration("latency", time.Since(start)),
		)
	})
}
