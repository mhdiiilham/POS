package api

import (
	"context"
	"net/http"
	"os"

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

	mux.HandleFunc("/api/login", s.Login).
		Methods(http.MethodPost)

	return mux
}

func (s *server) CORS(mux http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedHeaders: []string{"*"},
	}).Handler(mux)
}

func (s *server) HandlerLogging(mux http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, mux)
}
