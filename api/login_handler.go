package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type (
	LoginResponse struct {
		AccessToken    string    `json:"accessToken"`
		TokenType      string    `json:"tokenType"`
		TokenExpiresIn time.Time `json:"tokenExpiresIn"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Passwrod string `json:"password"`
	}
)

func (s *server) Login(w http.ResponseWriter, r *http.Request) {
	const ops = "api.service.Login"
	ctx := context.WithValue(r.Context(), logger.RequestIDKey, uuid.New().String())
	var req LoginRequest

	time.Sleep(35 * time.Second)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error(ctx, ops, "error decode request body: %v", err)
		UnknownErrorResponse(w, err)
		return
	}

	if req.Email == "" || req.Passwrod == "" {
		FailedResponse(w, user.ErrEmptyEmailAndPassword, http.StatusBadRequest)
		return
	}

	accessToken, err := s.userService.Login(ctx, req.Email, req.Passwrod)
	if err != nil {
		if errors.Is(err, user.ErrInvalidEmailAndPasword) {
			FailedResponse(w, err, http.StatusBadRequest)
			return
		}

		logger.Error(ctx, ops, "unknown: %v", err)
		UnknownErrorResponse(w, err)
		return
	}

	logger.Info(ctx, ops, "user %s login", req.Email)
	resp := LoginResponse{
		AccessToken:    accessToken,
		TokenType:      "Bearer",
		TokenExpiresIn: time.Now().Add(12 * time.Hour),
	}
	SuccessResponse(w, "login success", resp, http.StatusOK)
}
