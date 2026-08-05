package main

import (
	"context"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/middleware"

	"github.com/0x-ximon/portman/api/handlers"
)

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()

	cfg := Config{}
	err := cfg.Load(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	pool := cfg.pool
	defer pool.Close()

	// Goroutines
	go cfg.batcher.Start(ctx)

	// Middlewares
	mid := Middleware{}
	chain := mid.NewChain(
		mid.Logging,
		mid.Auth,
		chi.Logger,
	)

	addr := cfg.addr
	server := http.Server{
		Addr:    addr,
		Handler: chain(mux),
	}

	// Handlers
	deps := &handlers.Dependencies{
		DB:      pool,
		Web3:    cfg.web3,
		Mailer:  cfg.mailer,
		Cacher:  cfg.cacher,
		Batcher: cfg.batcher,
	}

	auth := deps.NewAuthHandler()
	mux.HandleFunc("POST /auth/initiate", auth.Initiate)
	mux.HandleFunc("POST /auth/validate", auth.Validate)
	mux.HandleFunc("POST /auth/exchange", auth.Exchange)

	users := deps.NewUsersHandler()
	mux.HandleFunc("GET /users", users.List)
	mux.HandleFunc("POST /users", users.Create)
	mux.HandleFunc("GET /users/{id}", users.Get)
	mux.HandleFunc("POST /users/fund", users.Fund)
	mux.HandleFunc("DELETE /users/{id}", users.Delete)

	tickers := deps.NewTickerHandler()
	mux.HandleFunc("GET /tickers", tickers.List)
	mux.HandleFunc("POST /tickers", tickers.Create)
	mux.HandleFunc("GET /tickers/{id}", tickers.Get)
	mux.HandleFunc("DELETE /tickers/{id}", tickers.Delete)
	mux.HandleFunc("GET /tickers/tick", tickers.Tick)

	orders := deps.NewOrderHandler()
	mux.HandleFunc("GET /orders", orders.List)
	mux.HandleFunc("POST /orders", orders.Create)
	mux.HandleFunc("GET /orders/{id}", orders.Get)

	log.Printf("Portman API listening on %s\n", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
