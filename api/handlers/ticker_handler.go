package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TickerHandler struct {
	DbConn *pgxpool.Pool
}

func (h *TickerHandler) GetTicker(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Payload{
			Message: "invalid id",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	ticker, err := repo.GetTicker(ctx, int32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "ticker not found",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Payload{
		Message: "ticker retrieved",
		Data:    ticker,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *TickerHandler) CreateTicker(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	var params repositories.CreateTickerParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Payload{
			Message: "invalid params",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	ticker, err := repo.CreateTicker(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create ticker",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Payload{
		Message: "ticker created",
		Data:    ticker,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *TickerHandler) ListTickers(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	tickers, err := repo.ListTickers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not list tickers",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Payload{
		Message: "tickers retrieved",
		Data:    tickers,
	}

	json.NewEncoder(w).Encode(results)
}

func (h *TickerHandler) DeleteTicker(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		results := Payload{
			Message: "invalid id",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(results)
		return
	}

	err = repo.DeleteTicker(ctx, int32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		results := Payload{
			Message: "could not delete ticker",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(results)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Payload{
		Message: "ticker deleted",
		Data:    nil,
	}

	json.NewEncoder(w).Encode(results)
}

func (h *TickerHandler) Tick(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := context.Background()

	_, _ = repo, ctx

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if os.Getenv("ENV") == "dev" {
				return true
			}

			origin, allowedOrigin := r.Header.Get("Origin"), os.Getenv("ALLOWED_ORIGIN") // "https://agence.ximon.dev"
			if origin == allowedOrigin {
				return true
			}

			log.Printf("Blocked unauthorized WebSocket connection attempt from: %s", origin)
			return false
		},
	}

	_ = upgrader
}
