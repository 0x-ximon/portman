package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	db     *pgxpool.Pool
	mailer *services.MailService
	cacher *services.CacheService
}

func (h *AuthHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	type Credentials struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	var params Credentials
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode login body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	user, err := repo.FindUserByEmail(ctx, params.EmailAddress)
	if err != nil {
		logger.Warn("failed to get user", "email", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "invalid credentials")
		return
	}

	if !services.ValidateHash(params.Password, user.Password) {
		logger.Warn("invalid login attempt", "email", params.EmailAddress)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "invalid credentials")
		return
	}

	otp, err := services.GenerateOTP(6)
	if err != nil {
		logger.Error("otp generation failed", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	if err := h.mailer.SendOTP(params.EmailAddress, otp); err != nil {
		logger.Error("email delivery failed", "user_id", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "Failed to send verification code")
		return
	}

	if err := h.cacher.StoreOTP(ctx, params.EmailAddress, otp); err != nil {
		logger.Error("cache storage failed", "email_address", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("otp sent successfully", "email_address", params.EmailAddress)
	SendSuccess(w, nil)
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	type Credentials struct {
		EmailAddress string `json:"email_address"`
		OTP          string `json:"otp"`
	}

	var params Credentials
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode login body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	otp, err := h.cacher.RetrieveOTP(ctx, params.EmailAddress)
	if err != nil {
		logger.Error("cache storage retrieval failed", "email_address", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	if otp != params.OTP {
		logger.Warn("invalid otp verification attempt", "email_address", params.EmailAddress, "otp", params.OTP)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "Invalid OTP")
		return
	}

	user, err := repo.FindUserByEmail(ctx, params.EmailAddress)
	if err != nil {
		logger.Warn("failed to get user", "email", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "invalid credentials")
		return
	}

	jwt, err := services.GenerateJWT(&user)
	if err != nil {
		logger.Error("failed to generate jwt", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	if err := h.cacher.DeleteOTP(ctx, params.EmailAddress); err != nil {
		logger.Error("failed to delete otp", "email_address", params.EmailAddress, "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("otp verified successfully", "email_address", params.EmailAddress)
	SendSuccess(w, jwt)
}

func (h *AuthHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	// TODO: Switch to asymmetric encryption and validate signature
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		logger.Warn("api-key not set in headers")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Missing X-API-KEY header")
		return
	}

	user, err := repo.FindUserByApiKey(ctx, &apiKey)
	if err == nil {
		logger.Warn("failed to get user", "error", err)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "invalid api-key")
		return
	}

	jwt, err := services.GenerateJWT(&user)
	if err != nil {
		logger.Error("failed to generate jwt", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("api-key exchanged successfully", "email_address", user.EmailAddress)
	SendSuccess(w, jwt)
}
