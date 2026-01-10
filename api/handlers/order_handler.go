package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/proto"
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
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}
	userID := claims.ID

	orderID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Payload{
			Message: "invalid id",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	params := repositories.GetOrderParams{ID: orderID, UserID: userID}
	order, err := repo.GetOrder(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "order not found",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Payload{
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
		result := Payload{
			Message: "invalid params",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {

		w.WriteHeader(http.StatusUnauthorized)
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	user, err := repo.GetUser(ctx, claims.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "user not found",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	// TODO: Implement user balance check
	_ = user

	order, err := repo.CreateOrder(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create order",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	serviceClient := proto.NewOrdersServiceClient(h.CoreConn)

	_, err = serviceClient.SubmitOrder(ctx, &proto.SubmitOrderRequest{
		Id: order.ID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not send order to core",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Payload{
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
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	userID := claims.ID

	orders, err := repo.ListOrders(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not list orders",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Payload{
		Message: "orders retrieved",
		Data:    orders,
	}

	json.NewEncoder(w).Encode(results)
}
