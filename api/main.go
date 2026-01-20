package main

import (
	"context"
	"log"
	"net/http"

	"github.com/0x-ximon/portman/api/handlers"
	"github.com/0x-ximon/portman/api/services"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
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

	cfg := Config{}
	err := cfg.Load(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	dc := cfg.dbConn
	defer dc.Close(ctx)

	cc := cfg.coreConn
	defer cc.Close()

	nc := cfg.natsConn
	defer nc.Close()

	addr := cfg.addr
	chain := services.NewChain(
		services.ContentType,
		services.Auth,

		middleware.Logger,
		middleware.Heartbeat("/health"),
	)

	server := http.Server{
		Addr:    addr,
		Handler: chain(mux),
	}

	auth := &handlers.AuthHandler{DbConn: dc}
	mux.HandleFunc("POST /auth/initiate", auth.Initiatiate)
	mux.HandleFunc("POST /auth/validate", auth.Validate)

	users := &handlers.UsersHandler{DbConn: dc}
	mux.HandleFunc("GET /users", users.ListUsers)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /users/{id}", users.GetUser)
	mux.HandleFunc("DELETE /users/{id}", users.DeleteUser)

	tickers := &handlers.TickerHandler{DbConn: dc}
	mux.HandleFunc("GET /tickers", tickers.ListTickers)
	mux.HandleFunc("POST /tickers", tickers.CreateTicker)
	mux.HandleFunc("GET /tickers/{id}", tickers.GetTicker)
	mux.HandleFunc("DELETE /tickers/{id}", tickers.DeleteTicker)

	orders := &handlers.OrderHandler{DbConn: dc, CoreConn: cc}
	mux.HandleFunc("GET /orders", orders.ListOrders)
	mux.HandleFunc("POST /orders", orders.CreateOrder)
	mux.HandleFunc("GET /orders/{id}", orders.GetOrder)

	// TODO: Move Stream Consumption into its own microservice
	cons, err := cfg.GetConsumers(ctx)
	if err != nil {
		log.Println(err)
	} else {
		cons.ordersProcessed.Consume(orders.ProcessOrder)
	}

	log.Printf("Starting server on %s", addr)
	server.ListenAndServe()
}
