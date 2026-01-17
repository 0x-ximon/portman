package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/0x-ximon/portman/api/handlers"
	"github.com/0x-ximon/portman/api/services"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	mux := http.NewServeMux()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctx)

	coreConn, err := services.NewCoreService().Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer coreConn.Close()

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	defer nc.Close()

	chain := services.NewChain(
		services.ContentType,
		services.Auth,

		middleware.Logger,
		middleware.Heartbeat("/health"),
	)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3001"
	}

	addr := net.JoinHostPort(os.Getenv("HOST"), port)
	s := http.Server{
		Addr:    addr,
		Handler: chain(mux),
	}

	auth := &handlers.AuthHandler{Conn: conn}
	mux.HandleFunc("POST /auth/initiate", auth.Initiatiate)
	mux.HandleFunc("POST /auth/validate", auth.Validate)

	users := &handlers.UsersHandler{Conn: conn}
	mux.HandleFunc("GET /users", users.ListUsers)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /users/{id}", users.GetUser)
	mux.HandleFunc("DELETE /users/{id}", users.DeleteUser)

	tickers := &handlers.TickerHandler{Conn: conn}
	mux.HandleFunc("GET /tickers", tickers.ListTickers)
	mux.HandleFunc("POST /tickers", tickers.CreateTicker)
	mux.HandleFunc("GET /tickers/{id}", tickers.GetTicker)
	mux.HandleFunc("DELETE /tickers/{id}", tickers.DeleteTicker)

	orders := &handlers.OrderHandler{Conn: conn, CoreConn: coreConn}
	mux.HandleFunc("GET /orders", orders.ListOrders)
	mux.HandleFunc("POST /orders", orders.CreateOrder)
	mux.HandleFunc("GET /orders/{id}", orders.GetOrder)
	nc.Subscribe("orders.processed", orders.ProcessOrder)

	log.Printf("Starting server on %s", addr)
	s.ListenAndServe()
}
