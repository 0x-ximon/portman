package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
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

// func (h *TickerHandler) Tick(w http.ResponseWriter, r *http.Request) {
// 	repo := repositories.New(h.DbConn)
// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			if os.Getenv("ENV") == "dev" {
// 				return true
// 			}

// 			origin, allowedOrigin := r.Header.Get("Origin"), os.Getenv("ALLOWED_ORIGIN") // "https://agence.ximon.dev"
// 			if origin == allowedOrigin {
// 				return true
// 			}

// 			log.Printf("Blocked unauthorized WebSocket connection attempt from: %s", origin)
// 			return false
// 		},
// 	}

// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		result := Payload{
// 			Message: "websocket upgrade error",
// 			Error:   err.Error(),
// 		}

// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}
// 	defer ws.Close()

// 	symbolChan := make(chan string, 1)
// 	go func() {
// 		for {
// 			_, msg, err := ws.ReadMessage()
// 			if err != nil {
// 				cancel()
// 				return
// 			}

// 			symbolChan <- string(msg)
// 		}
// 	}()

// 	var sub *nats.Subscription
// 	msgChan := make(chan *nats.Msg, 64)

// 	for {
// 		select {
// 		case symbol := <-symbolChan:
// 			if sub != nil {
// 				sub.Unsubscribe()
// 			}

// 			// TODO: Properly handle Ticker WebSocket Subscription Errors
// 			ticker, err := repo.FindTickerBySymbol(ctx, symbol)
// 			if err != nil {
// 				continue
// 			}

// 			if ticker.Status != repositories.TickerStatusOPEN {
// 				continue
// 			}

// 			sub, err = h.NatsConn.ChanSubscribe(fmt.Sprintf("ticks.%s", symbol), msgChan)
// 			if err != nil {
// 				continue
// 			}

// 		case msg := <-msgChan:
// 			err := ws.WriteMessage(websocket.TextMessage, msg.Data)
// 			if err != nil {
// 				return
// 			}

// 		case <-ctx.Done():
// 			sub.Unsubscribe()
// 			return
// 		}
// 	}
// }
