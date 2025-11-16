package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	Conn *pgx.Conn
}

func (h *AuthHandler) Initiatiate(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	var params Credentials
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

	user, err := repo.FindUserByEmail(ctx, params.EmailAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "user not found",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if !services.ValidateHash(params.Password, user.Password) {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "invalid credentials",
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

	if err := cacher.StoreOTP(ctx, user.ID, otp); err != nil {
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
		Message: "otp sent",
	}

	json.NewEncoder(w).Encode(result)
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	repo := repositories.New(h.Conn)
	ctx := r.Context()

	var params Credentials
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

	user, err := repo.FindUserByEmail(ctx, params.EmailAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "user not found",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	otp, err := cacher.RetrieveOTP(ctx, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not get otp",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if otp != params.OTP {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "invalid otp",
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	token, err := services.GenerateJWT(user.EmailAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not generate token",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	if err := cacher.DeleteOTP(ctx, user.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result := Result{
			Message: "could not delete otp",
			Error:   err,
		}

		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(http.StatusOK)
	result := Result{
		Message: "otp validated",
		Data:    token,
	}

	json.NewEncoder(w).Encode(result)
}
