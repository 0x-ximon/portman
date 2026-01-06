package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/0x-ximon/portman/api/handlers"
	"github.com/0x-ximon/portman/api/middlewares"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	router := http.NewServeMux()
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

	chain := middlewares.NewChain(
		middlewares.ContentType,
		middlewares.Auth,
		middlewares.Logger,
	)

	auth := &handlers.AuthHandler{Conn: conn}
	router.HandleFunc("POST /auth/initiate", auth.Initiatiate)
	router.HandleFunc("POST /auth/validate", auth.Validate)

	users := &handlers.UsersHandler{Conn: conn}
	router.HandleFunc("GET /users", users.ListUsers)
	router.HandleFunc("POST /users", users.CreateUser)
	router.HandleFunc("GET /users/{id}", users.GetUser)
	router.HandleFunc("DELETE /users/{id}", users.DeleteUser)

	tickers := &handlers.TickerHandler{Conn: conn}
	router.HandleFunc("GET /tickers", tickers.ListTickers)
	router.HandleFunc("POST /tickers", tickers.CreateTicker)
	router.HandleFunc("GET /tickers/{id}", tickers.GetTicker)
	router.HandleFunc("DELETE /tickers/{id}", tickers.DeleteTicker)

	orders := &handlers.OrderHandler{Conn: conn, CoreConn: coreConn}
	router.HandleFunc("GET /orders", orders.ListOrders)
	router.HandleFunc("POST /orders", orders.CreateOrder)
	router.HandleFunc("GET /orders/{id}", orders.GetOrder)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3001"
	}

	addr := net.JoinHostPort(os.Getenv("HOST"), port)
	s := http.Server{
		Addr:    addr,
		Handler: chain(router),
	}

	log.Printf("Starting server on %s", addr)
	s.ListenAndServe()
}
