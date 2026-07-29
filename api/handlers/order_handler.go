package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	db      *pgxpool.Pool
	mailer  *services.MailService
	batcher *services.BatchService
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		logger.Warn("failed to parse order id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	params := repositories.GetOrderParams{ID: id, UserID: claims.ID}
	order, err := repo.GetOrder(ctx, params)
	if err != nil {
		logger.Warn("failed to get order", "order_id", id, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "order not found")
		return
	}

	logger.Info("order retrieved successfully", "order_id", id, "user_id", claims.ID)
	SendSuccess(w, order)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	var params repositories.CreateOrderParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode create order body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	user, err := repo.GetUser(ctx, claims.ID)
	if err != nil {
		logger.Warn("failed to get user", "user_id", claims.ID, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "user not found")
		return
	}

	// TODO: Enable this and make this work on all currencies
	// if user.FreeBalance.LessThan(params.Price.Mul(params.Quantity)) {
	// 	logger.Warn("insufficient balance", "user_id", user.ID)
	// 	SendFailure(w, http.StatusNotFound, codeNotFound, "insufficient balance")
	// 	return
	// }

	order, err := repo.CreateOrder(ctx, params)
	if err != nil {
		logger.Error("failed to create order", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	h.batcher.CreateOrder <- order
	logger.Info("order created successfully", "order_id", order.ID, "user_id", user.ID)
	SendSuccess(w, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	orders, err := repo.ListOrders(ctx, claims.ID)
	if err != nil {
		logger.Warn("failed to get orders", "user_id", claims.ID, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "orders not found")
		return
	}

	logger.Info("orders retrieved successfully", "user_id", claims.ID)
	SendSuccess(w, orders)
}

// func (h *OrderHandler) Process(msg jetstream.Msg) {
// 	repo := repositories.New(h.DbConn)
// 	ctx := context.Background()

// 	type Params struct {
// 		ID     int64
// 		Status repositories.OrderStatus
// 	}

// 	r := bytes.NewReader(msg.Data())
// 	var params []Params

// 	err := json.NewDecoder(r).Decode(&params)
// 	if err != nil {
// 		msg.Nak()
// 		return
// 	}

// 	for _, param := range params {
// 		args := repositories.UpdateOrderParams{
// 			ID:     param.ID,
// 			Status: param.Status,
// 		}

// 		err = repo.UpdateOrder(ctx, args)
// 		if err != nil {
// 			msg.Nak()
// 			return
// 		}
// 	}

// 	msg.Ack()
// }
