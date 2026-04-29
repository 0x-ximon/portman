package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TickerHandler struct {
	db *pgxpool.Pool
}

func (h *TickerHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		logger.Warn("failed to parse ticker id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	ticker, err := repo.GetTicker(ctx, int32(id))
	if err != nil {
		logger.Warn("failed to get ticker", "ticker_id", id, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "ticker  not found")
		return
	}

	logger.Info("ticker retrieved successfully", "ticker_id", ticker.ID)
	SendSuccess(w, ticker)
}

func (h *TickerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	var params repositories.CreateTickerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode create ticker body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	// TODO: Create the Ticker in the core

	ticker, err := repo.CreateTicker(ctx, params)
	if err != nil {
		logger.Error("failed to create ticker", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("ticker created successfully", "ticker_id", ticker.ID)
	SendSuccess(w, ticker)
}

func (h *TickerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		logger.Warn("failed to parse ticker id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	if claims.Role != repositories.RoleADMINISTRATOR {
		logger.Warn("user is not an admin", "user_id", claims.ID, "role", claims.Role)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user not authorized to list users")
		return
	}

	err = repo.DeleteTicker(ctx, int32(id))
	if err != nil {
		logger.Error("database error during ticker deletion", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("ticker deleted successfully", "ticker_id", id)
	SendSuccess(w, nil)
}

func (h *TickerHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	tickers, err := repo.ListTickers(ctx)
	if err != nil {
		logger.Error("database error during users lookup", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("tickers retrieved successfully")
	SendSuccess(w, tickers)
}

func (h *TickerHandler) Tick(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if os.Getenv("ENV") == "dev" {
				return true
			}

			origin, allowedOrigin := r.Header.Get("Origin"), os.Getenv("ALLOWED_ORIGIN") // "https://agence.ximon.dev"
			if origin == allowedOrigin {
				return true
			}

			logger.Warn("attempted unauthorized websocket connection", "origin", origin)
			return false
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warn("failed to upgrade connection to websocket", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "Something went wrong")
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

	type Msg struct {
		Ask  int32 `json:"ask"`
		Bid  int32 `json:"bid"`
		Last int32 `json:"last"`
	}

	msgChan := make(chan *Msg, 64)
	for {
		select {
		case symbol := <-symbolChan:
			_, err := repo.FindTickerBySymbol(ctx, symbol)
			if err != nil {
				return
			}

			go func(startPrice int32) {
				currentPrice := startPrice
				t := time.NewTicker(2 * time.Second)
				defer t.Stop()

				for {
					select {

					case <-ctx.Done():
						return

					case <-t.C:
						change := int32(rand.Intn(5) - 2)
						currentPrice += change
						if currentPrice < 1 {
							currentPrice = 1
						}

						msgChan <- &Msg{
							Last: currentPrice,
							Ask:  currentPrice + 1,
							Bid:  currentPrice - 1,
						}
					}
				}
			}(100)

		case msg := <-msgChan:
			var buffer bytes.Buffer
			json.NewEncoder(&buffer).Encode(msg)
			if err := ws.WriteMessage(websocket.TextMessage, buffer.Bytes()); err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// go func() {
// 	for {
// 		select {
//
// 		case <-ctx.Done():
// 			return
//
// 		default:
// 			msg := Msg{Ask: 120, Bid: 100, Last: 110}
// 			time.Sleep(time.Second * 3)
// 			msgChan <- &msg
// 		}
// 	}
// }()

// 	for {
// 		select {
// 		case symbol := <-symbolChan:
// 			// TODO: Properly handle Ticker WebSocket Subscription Errors
// 			ticker, err := repo.FindTickerBySymbol(ctx, symbol)
// 			if err != nil {
// 				return
// 			}
//
// 			if ticker.Status != repositories.TickerStatusOPEN {
// 				continue
// 			}
//
// 		case msg := <-msgChan:
// 			var buffer bytes.Buffer
// 			if err := json.NewEncoder(&buffer).Encode(msg); err != nil {
// 				logger.Warn("failed to encode message", "error", err)
// 				continue
// 			}
//
// 			err := ws.WriteMessage(websocket.TextMessage, buffer.Bytes())
// 			if err != nil {
// 				return
// 			}
//
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
