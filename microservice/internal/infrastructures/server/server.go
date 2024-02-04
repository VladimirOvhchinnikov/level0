package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	router *chi.Mux
	server *http.Server
}

func NewGoChi(logger *zap.Logger, router *chi.Mux) *Server {

	return &Server{
		logger: logger,
		router: router,
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

func (srv *Server) StartServer(port string) error {

	if port == "" {
		srv.logger.Info("Port is empty. Initializing standard values")
		srv.server.Addr = ":8080"
	}

	if srv.router == nil {
		srv.logger.Error("Routing is not initialized")
		return errors.New("Routing is not initialized")
	}

	srv.logger.Info("Starting server on port", zap.String("port", port))
	if err := srv.server.ListenAndServe(); err != http.ErrServerClosed {
		srv.logger.Error("Error when starting the server:", zap.Error(err))
		return err
	}

	return nil
}

func (gc *Server) Shutdown() error {

	if gc.server != nil {
		if err := gc.server.Shutdown(context.Background()); err != nil {
			gc.logger.Error("Error when stopping the server:", zap.Error(err))
			return err
		}
	}

	return nil
}
