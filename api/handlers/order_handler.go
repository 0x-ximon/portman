package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	db     *pgxpool.Pool
	mailer *services.MailService
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		logger.Warn("failed to parse order id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	userID := uuid.Nil
	user, ok := ctx.Value(services.UserKey{}).(repositories.User)
	if ok {
		userID = user.ID
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(services.Claims)
	if !ok || userID == uuid.Nil {
		logger.Warn("attempted order retrieval without authorization", "user_id", claims.ID)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user not authorized to retrieve order")
		return
	}
	userID = claims.ID

	params := repositories.GetOrderParams{ID: orderID, UserID: userID}
	order, err := repo.GetOrder(ctx, params)
	if err != nil {
		logger.Warn("failed to get order", "order_id", orderID, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "order not found")
		return
	}

	logger.Info("order retrieved successfully", "order_id", orderID, "user_id", userID)
	SendSuccess(w, order)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	userID := uuid.Nil
	user, ok := ctx.Value(services.UserKey{}).(repositories.User)
	if ok {
		userID = user.ID
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(services.Claims)
	if !ok || userID == uuid.Nil {
		logger.Warn("attempted order retrieval without authorization", "user_id", claims.ID, "role", claims.Role)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user not authorized to retrieve order")
		return
	}
	userID = claims.ID

	var params repositories.CreateOrderParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode create order body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	user, err := repo.GetUser(ctx, userID)
	if err != nil {
		logger.Warn("failed to get user", "user_id", userID, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "user not found")
		return
	}

	// TODO: Make this work on all currencies
	if user.FreeBalance.LessThan(params.Price.Mul(params.Quantity)) {
		logger.Warn("insufficient balance", "user_id", userID)
		SendFailure(w, http.StatusNotFound, codeNotFound, "insufficient balance")
		return
	}

	order, err := repo.CreateOrder(ctx, params)
	if err != nil {
		logger.Error("failed to create order", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	// TODO: Send order to core

	logger.Info("order created successfully", "order_id", order.ID, "user_id", user.ID)
	SendSuccess(w, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	userID := uuid.Nil
	user, ok := ctx.Value(services.UserKey{}).(repositories.User)
	if ok {
		userID = user.ID
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(services.Claims)
	if !ok || userID == uuid.Nil {
		logger.Warn("attempted orders retrieval without authorization", "user_id", claims.ID, "role", claims.Role)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user not authorized to retrieve order")
		return
	}
	userID = claims.ID

	orders, err := repo.ListOrders(ctx, userID)
	if err != nil {
		logger.Warn("failed to get orders", "user_id", userID, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "orders not found")
		return
	}

	logger.Info("orders retrieved successfully", "user_id", userID)
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
