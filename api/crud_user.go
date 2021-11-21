package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type (
	CreateUserRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	CreateUserResponse struct {
		User user.User `json:"user"`
	}
)

func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {
	const ops = "api.server.CreateUser"
	ctx := context.WithValue(r.Context(), logger.RequestIDKey, uuid.New().String())
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		UnknownErrorResponse(w, err)
		return
	}

	userCredential, ok := r.Context().Value("user-credentials").(TokenPayload)
	if !ok {
		err := errors.New("failed to cast contect value of user-credentials to type TokenPayload")
		logger.Error(ctx, ops, err.Error())
		UnknownErrorResponse(w, err)
		return
	}

	entity := user.User{
		MerchantID: userCredential.MerchantID,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   &req.LastName,
		Password:   req.Password,
	}

	uID, err := s.userService.CreateUser(ctx, entity)
	if err != nil {
		logger.Info(ctx, ops, "err: %v", err)
		if errors.Is(err, user.ErrEmailNotUnique) {
			FailedResponse(w, err, http.StatusConflict)
			return
		}

		if errors.Is(err, user.ErrInvalidCreateParameters) {
			FailedResponse(w, err, http.StatusBadRequest)
			return
		}

		UnknownErrorResponse(w, err)
		return
	}

	entity.ID = uID
	resp := CreateUserResponse{
		User: entity,
	}

	logger.Info(ctx, ops, "success created new user")
	SuccessResponse(w, "success", resp, http.StatusCreated)
}
