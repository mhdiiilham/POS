package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	GetUsersResponse struct {
		Users     []user.User `json:"users"`
		Page      int         `json:"page"`
		TotalData int         `json:"totalData"`
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

func (s *server) GetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		page      int = 1
		limit     int
		lastID    int
		totalData int
		err       error
		users     []user.User = []user.User{}
	)

	const ops = "api.server.GetUsers"
	ctx := context.WithValue(r.Context(), logger.RequestIDKey, uuid.New().String())
	userCredentials := r.Context().Value("user-credentials").(TokenPayload)
	limitQuery := r.URL.Query().Get("limit")
	lastIDQuery := r.URL.Query().Get("lastID")
	pageQuery := r.URL.Query().Get("page")

	logger.Info(ctx, ops, "start handling GetUsers")
	limit, err = strconv.Atoi(limitQuery)
	if err != nil && limitQuery != "" {
		FailedResponse(w, errors.New("invalid limit"), http.StatusBadRequest)
		return
	}

	lastID, err = strconv.Atoi(lastIDQuery)
	if err != nil && lastIDQuery != "" {
		FailedResponse(w, errors.New("invalid lastid"), http.StatusBadRequest)
		return
	}

	page, err = strconv.Atoi(pageQuery)
	if err != nil && pageQuery != "" {
		FailedResponse(w, errors.New("invalid page"), http.StatusBadRequest)
		return
	}

	if limitQuery == "" {
		limit = 10
	}

	if pageQuery == "" {
		page = 1
	}

	users, totalData, err = s.userService.GetUsers(ctx, userCredentials.MerchantID, lastID, limit)
	if err != nil {
		logger.Error(ctx, ops, "unexpected error %v", err)
		UnknownErrorResponse(w, err)
		return
	}

	SuccessResponse(w, "data found", GetUsersResponse{
		Page:      page,
		Users:     users,
		TotalData: totalData,
	}, http.StatusOK)
}

func (s *server) RemoveUser(w http.ResponseWriter, r *http.Request) {
	const ops = "api.service.RemoveUser"
	ctx := context.WithValue(r.Context(), logger.RequestIDKey, uuid.New().String())

	vars := mux.Vars(r)
	userIDParam := vars["userId"]
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		FailedResponse(w, errors.New("invalid user id"), http.StatusBadRequest)
		return
	}

	err = s.userService.DeleteUser(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			FailedResponse(w, errors.New("invalid user id"), http.StatusBadRequest)
			return
		}

		logger.Error(ctx, ops, "unkown error: %v", err.Error())
		UnknownErrorResponse(w, err)
		return
	}

	SuccessResponse(w, fmt.Sprintf("success delete user with id %d", userID), nil, http.StatusOK)
}

func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {
	const ops = "api.service.RemoveUser"
	ctx := context.WithValue(r.Context(), logger.RequestIDKey, uuid.New().String())

	vars := mux.Vars(r)
	userIDParam := vars["userId"]
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		FailedResponse(w, errors.New("invalid user id"), http.StatusBadRequest)
		return
	}

	entity, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			FailedResponse(w, user.ErrUserNotFound, http.StatusBadRequest)
			return
		}

		logger.Error(ctx, ops, "unkown error: %v", err.Error())
		UnknownErrorResponse(w, err)
		return
	}

	SuccessResponse(w, "data found", entity, http.StatusOK)
}
