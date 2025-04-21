package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/delivery/http/router"
	"github.com/WebChads/AccountService/internal/usecase"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type Server struct {
	server *http.Server
}

// Server constructor
func NewServer(handler http.Handler, address string) *Server {
	srv := &http.Server{
		Addr:           address,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{server: srv}
}

func NewDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitRouter(config *config.ServerConfig, logger *slog.Logger, db *sql.DB) http.Handler {
	rout := chi.NewRouter()
	http.Handle("/", rout)

	repos := usecase.NewRepositories(db)

	// Add all routers here
	accountUsecase := usecase.NewAccountUsecase(repos.Account)
	accountRouter := router.NewAccountRouter(rout, config, logger, accountUsecase)
	// ...

	// Configure routers
	router.ConfigureAccountRouter(accountRouter)
	// ...

	return rout
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
