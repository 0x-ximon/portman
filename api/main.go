package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/0x-ximon/portman/api/handlers"
	"github.com/0x-ximon/portman/api/middleware"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	host := os.Getenv("HOST")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3001"
	}

	addr := net.JoinHostPort(host, port)
	router := http.NewServeMux()

	router.HandleFunc("GET /tickers", handlers.GetTickers)
	router.HandleFunc("POST /tickers", handlers.CreateTicker)
	router.HandleFunc("GET /tickers/{id}", handlers.GetTicker)
	router.HandleFunc("DELETE /tickers/{id}", handlers.DeleteTicker)

	s := http.Server{
		Addr:    addr,
		Handler: middleware.Logger(router),
	}

	log.Printf("Starting server on %s", addr)
	s.ListenAndServe()
}
