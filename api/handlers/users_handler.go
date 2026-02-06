package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
)

type UsersHandler struct {
	DbConn *pgx.Conn
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	user, ok := r.Context().Value(services.UserKey{}).(*repositories.User)
	if ok {
		w.WriteHeader(http.StatusOK)
		result := Payload{
			Message: "user retrieved",
			Data:    user,
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

	w.WriteHeader(http.StatusOK)
	result := Payload{
		Message: "user retrieved",
		Data:    user,
	}

	json.NewEncoder(w).Encode(result)
	return

}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	var params repositories.CreateUserParams
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

	encryptedPassword, err := services.HashPassword(params.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not hash password",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}
	params.Password = encryptedPassword

	user, err := repo.CreateUser(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create user",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	otp, err := services.GenerateOTP(6)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not generate otp",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	mailer, err := services.NewMailService()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create mailer",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	cacher, err := services.NewCacheService()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not create cacher",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if err := mailer.SendOTP(user.EmailAddress, otp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not send otp",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if err := cacher.StoreOTP(ctx, user.ID, otp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Payload{
			Message: "could not set otp",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Payload{
		Message: "user created",
		Data:    user,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.DbConn)
	ctx := r.Context()

	id, ok := services.GetIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		result := Payload{
			Message: "unauthorized",
			Error:   "bearer token not found",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	err := repo.DeleteUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		results := Payload{
			Message: "could not delete user",
			Error:   err.Error(),
		}

		json.NewEncoder(w).Encode(results)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Payload{
		Message: "user deleted",
	}

	json.NewEncoder(w).Encode(results)
}
