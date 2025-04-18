package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/WebChads/AccountService/internal/delivery/http/handler"
	"github.com/go-chi/chi"
)

type Server struct {
	server *http.Server
}

// Server constructor
func NewServer(port int) *Server {
	router := InitRouter()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{server: srv}
}

func InitRouter() *chi.Mux {
	router := chi.NewRouter()

	account := handler.NewAccountHandler()
	router.HandleFunc("/api/v1/account/create-account", account.CreateAccountHandler)
	router.HandleFunc("/api/v1/account/update-account", account.UpdateAccountHandler)

	return router
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
