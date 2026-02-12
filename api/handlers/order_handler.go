package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x-ximon/portman/api/proto"
	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
)

var sideMap = map[repositories.OrderSide]proto.Side{
	repositories.OrderSideBUY:  proto.Side_SIDE_BUY,
	repositories.OrderSideSELL: proto.Side_SIDE_SELL,
}

var typeMap = map[repositories.OrderType]proto.Type{
	repositories.OrderTypeLIMIT:  proto.Type_TYPE_LIMIT,
	repositories.OrderTypeMARKET: proto.Type_TYPE_MARKET,
}

var statusMap = map[repositories.OrderStatus]proto.Status{
	repositories.OrderStatusPENDING:   proto.Status_STATUS_PENDING,
	repositories.OrderStatusCANCELLED: proto.Status_STATUS_CANCELLED,
	repositories.OrderStatusFULFILLED: proto.Status_STATUS_FULFILLED,
}

type OrderHandler struct {
	DbConn   *pgxpool.Pool
	CoreConn *grpc.ClientConn
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	userID, ok := services.GetIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

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
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	userID, ok := services.GetIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

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

	user, err := repo.GetUser(ctx, userID)
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

	request := proto.SubmitOrderRequest{
		Order: &proto.Order{
			Id:       order.ID,
			Side:     sideMap[order.Side],
			Type:     typeMap[order.Type],
			Status:   statusMap[order.Status],
			Price:    order.Price.String(),
			Quantity: order.Quantity.String(),
		},
		Symbol: order.TickerSymbol,
	}

	core := proto.NewOrdersServiceClient(h.CoreConn)
	response, err := core.SubmitOrder(ctx, &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not send order to core",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if response.Result != proto.Result_RESULT_SUCCESS {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "order not exectued",
			Error:   "core rejected the order",
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
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	userID, ok := services.GetIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

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

func (h *OrderHandler) ProcessOrder(msg jetstream.Msg) {
	repo := repositories.New(h.DbConn)
	ctx := context.Background()

	type Params struct {
		ID     int64
		Status repositories.OrderStatus
	}

	r := bytes.NewReader(msg.Data())
	var params []Params

	err := json.NewDecoder(r).Decode(&params)
	if err != nil {
		msg.Nak()
		return
	}

	for _, param := range params {
		args := repositories.UpdateOrderParams{
			ID:     param.ID,
			Status: param.Status,
		}

		err = repo.UpdateOrder(ctx, args)
		if err != nil {
			msg.Nak()
			return
		}
	}

	msg.Ack()
}
