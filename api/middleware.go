package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/mhdiiilham/POS/pkg/logger"
)

func (s *server) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const ops = "api.server.authorization"
		authorizationHeader := r.Header.Get("Authorization")

		if !strings.Contains(authorizationHeader, "Bearer") {
			err := errors.New("unauthorized")
			FailedResponse(w, err, http.StatusUnauthorized)
			return
		}

		signedToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		claims, err := s.tokenSigner.Extract(r.Context(), signedToken)
		if err != nil {
			UnknownErrorResponse(w, err)
			return
		}

		userID, userIDCastErr := claims["userID"].(float64)
		if !userIDCastErr {
			logger.Error(context.Background(), ops, "error casting userID to float64")
			UnknownErrorResponse(w, errors.New("cannot cast claims[userID] to float64"))
			return
		}
		merchantID, merchantIDCastErr := claims["merchantID"].(float64)
		if !merchantIDCastErr {
			logger.Error(context.Background(), ops, "error casting merchantID to float64")
			UnknownErrorResponse(w, errors.New("cannot cast claims[merchantID] to float64"))
			return
		}
		userEmail, userEmailCastErr := claims["email"].(string)
		if !userEmailCastErr {
			logger.Error(context.Background(), ops, "error casting usermail to string")
			UnknownErrorResponse(w, errors.New("cannot cast claims[email] to string"))
			return
		}

		data := TokenPayload{
			UserID:     int(userID),
			MerchantID: int(merchantID),
			Email:      userEmail,
		}

		ctx := context.WithValue(r.Context(), "user-credentials", data)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
