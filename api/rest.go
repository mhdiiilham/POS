package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
	"github.com/rs/cors"
)

type server struct {
	userService user.Service
}

func NewServer(userService user.Service) *server {
	return &server{userService: userService}
}

func (s *server) Routes(ctx context.Context) http.Handler {
	const ops = "api.Routes"

	logger.Info(ctx, ops, "initializing routing")
	mux := mux.NewRouter()

	mux.Use(s.APIMiddleware())
	mux.HandleFunc("/api/login", s.Login).
		Methods(http.MethodPost)

	JSON, _ := json.Marshal(Response{
		Code:    http.StatusRequestTimeout,
		Message: "request timeout",
		Data:    nil,
		Error:   http.ErrHandlerTimeout,
	})
	return http.TimeoutHandler(mux, 30*time.Second, string(JSON))
}

func (s *server) CORS(mux http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedHeaders: []string{"*"},
	}).Handler(mux)
}

func (s *server) HandlerLogging(mux http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, mux)
}

func (s *server) APIMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Date", time.Now().Format(time.RFC1123))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Connection", "keep-alive")

			next.ServeHTTP(w, r)
		})
	}
}
