package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/0x-ximon/portman/api/proto"
	"github.com/0x-ximon/portman/api/repositories"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type TickerHandler struct {
	DbConn   *pgxpool.Pool
	NatsConn *nats.Conn
	CoreConn *grpc.ClientConn
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

	request := proto.NewOrderBookRequest{
		// TODO: Modify these hardcoded values
		PricePrecision:    2,
		QuantityPrecision: 2,
		Symbol:            params.Symbol,
	}

	core := proto.NewOrdersServiceClient(h.CoreConn)
	response, err := core.NewOrderBook(ctx, &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create new order book in core",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if response.Result != proto.Result_RESULT_SUCCESS {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "ticker not created",
			Error:   "core rejected the ticker",
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
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

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

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "websocket upgrade error",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}
	defer ws.Close()

	symbolChan := make(chan string, 1)
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				cancel()
				return
			}

			symbolChan <- string(msg)
		}
	}()

	var sub *nats.Subscription
	msgChan := make(chan *nats.Msg, 64)

	for {
		select {
		case symbol := <-symbolChan:
			if sub != nil {
				sub.Unsubscribe()
			}

			// TODO: Properly handle Ticker WebSocket Subscription Errors
			ticker, err := repo.FindTickerBySymbol(ctx, symbol)
			if err != nil {
				continue
			}

			if ticker.Status != repositories.TickerStatusOPEN {
				continue
			}

			sub, err = h.NatsConn.ChanSubscribe(fmt.Sprintf("ticks.%s", symbol), msgChan)
			if err != nil {
				continue
			}

		case msg := <-msgChan:
			if err := ws.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
