package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	Conn     *pgx.Conn
	CoreConn *grpc.ClientConn
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Result{
			Message: "unauthorized",
			Error:   fmt.Errorf("bearer token not found"),
		}

		json.NewEncoder(w).Encode(result)
		return
	}
	userID := claims.ID

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Result{
			Message: "invalid id",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	params := repositories.GetOrderParams{ID: int32(orderID), BuyerID: userID}
	order, err := repo.GetOrder(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "order not found",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Result{
		Message: "order retrieved",
		Data:    order,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	var params repositories.CreateOrderParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Result{
			Message: "invalid params",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Result{
			Message: "unauthorized",
			Error:   fmt.Errorf("bearer token not found"),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	user, err := repo.GetUser(ctx, claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "user not found",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	// TODO: Implement user balance check
	_ = user

	order, err := repo.CreateOrder(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not create order",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	// TODO: Send Order To Core for Processing via gRPC

	w.WriteHeader(http.StatusOK)
	result := Result{
		Message: "order created",
		Data:    order,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Result{
			Message: "unauthorized",
			Error:   fmt.Errorf("bearer token not found"),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	userID := claims.ID

	orders, err := repo.ListOrders(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not list orders",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Result{
		Message: "orders retrieved",
		Data:    orders,
	}

	json.NewEncoder(w).Encode(results)
}
