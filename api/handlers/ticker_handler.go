package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/jackc/pgx/v5"
)

type TickerHandler struct {
	Conn *pgx.Conn
}

func (h *TickerHandler) GetTicker(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
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
	repo := repositories.New(h.Conn)
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
	repo := repositories.New(h.Conn)
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
	repo := repositories.New(h.Conn)
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
