package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/mhdiiilham/POS/pkg/logger"
)

type server struct {
	ip       string
	port     string
	listener net.Listener
}

func New(port string) (*server, error) {
	const scope = "server.New"

	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(context.Background(), scope, "failed to create listener on %s: %v", addr, err)
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	return &server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	const scope = "server.server.ServeHTTP"

	// Spawn a goroutine that listens for context closure. When the context is
	// closed, the server is stopped.
	errChan := make(chan error, 1)
	go func() {
		<-ctx.Done()

		shutDownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Info(ctx, scope, "server receive signal to shutdown")
		errChan <- srv.Shutdown(shutDownCtx)
	}()

	// Run the server. This will block until the provided context is closed.
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	return nil
}

func (s *server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		Handler: handler,
	})
}
