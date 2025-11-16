package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UsersHandler struct {
	Conn *pgx.Conn
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := Result{
			Message: "invalid id",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	user, err := repo.GetUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "user not found",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Result{
		Message: "user retrieved",
		Data:    user,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	var params repositories.CreateUserParams
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

	encryptedPassword, err := services.HashPassword(params.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not hash password",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}
	params.Password = encryptedPassword

	user, err := repo.CreateUser(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not create user",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	otp, err := services.GenerateOTP(6)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not generate otp",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	mailer, err := services.NewMailService()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not create mailer",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	cacher, err := services.NewCacheService()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not create cacher",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if err := mailer.SendOTP(user.EmailAddress, otp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not send otp",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if err := cacher.SetOTP(ctx, user.ID, otp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not set otp",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Result{
		Message: "user created",
		Data:    user,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	users, err := repo.ListUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not list users",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Result{
		Message: "users retrieved",
		Data:    users,
	}

	json.NewEncoder(w).Encode(results)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		results := Result{
			Message: "invalid id",
			Error:   err,
		}

		json.NewEncoder(w).Encode(results)
		return
	}

	err = repo.DeleteUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		results := Result{
			Message: "could not delete user",
			Error:   err,
		}

		json.NewEncoder(w).Encode(results)
		return
	}

	w.WriteHeader(http.StatusOK)
	results := Result{
		Message: "user deleted",
	}

	json.NewEncoder(w).Encode(results)
}
