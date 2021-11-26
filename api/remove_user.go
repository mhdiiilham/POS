package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
)

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
