package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/0x-ximon/portman/api/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersHandler struct {
	db *pgxpool.Pool
}

func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to parse user id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	if claims.ID != id {
		logger.Warn("claims id does not match user id", "user_id", id)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user id mismatch")
		return
	}

	user, err := repo.GetUser(ctx, id)
	if err != nil {
		logger.Warn("failed to get user", "user_id", id, "error", err)
		SendFailure(w, http.StatusNotFound, codeNotFound, "user not found")
		return
	}

	logger.Info("user retrieved successfully", "user_id", user.ID)
	SendSuccess(w, user)
}

func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	var params repositories.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode login body", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}

	encryptedPassword, err := services.HashPassword(params.Password)
	if err != nil {
		logger.Warn("failed to encrypt password", "error", err)
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "Invalid request format")
		return
	}
	params.Password = encryptedPassword

	user, err := repo.CreateUser(ctx, params)
	if err != nil {
		logger.Error("failed to create user", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("user created successfully", "user_id", user.ID)
	SendSuccess(w, user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		logger.Warn("failed to parse user id")
		SendFailure(w, http.StatusBadRequest, codeBadRequest, "invalid id")
		return
	}

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	if claims.ID != id {
		logger.Warn("claims id does not match user id", "user_id", id)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user id mismatch")
		return
	}

	err = repo.DeleteUser(ctx, id)
	if err != nil {
		logger.Error("failed to delete user", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("user deleted successfully", "user_id", id)
	SendSuccess(w, nil)
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repo := repositories.New(h.db)
	logger := services.GetLogger(ctx)

	claims, ok := r.Context().Value(services.ClaimsKey{}).(*services.Claims)
	if !ok {
		logger.Warn("failed to get claims from context")
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user claims not found")
		return
	}

	if claims.Role != repositories.UserRoleADMIN {
		logger.Warn("user is not an admin", "user_id", claims.ID, "role", claims.Role)
		SendFailure(w, http.StatusUnauthorized, codeUnauthorized, "user not authorized to list users")
		return
	}

	users, err := repo.ListUsers(ctx)
	if err != nil {
		logger.Error("database error during users lookup", "error", err)
		SendFailure(w, http.StatusInternalServerError, codeInternal, "An unexpected error occurred")
		return
	}

	logger.Info("users retrieved successfully")
	SendSuccess(w, users)
}
