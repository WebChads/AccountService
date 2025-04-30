package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/delivery/http/router"
	"github.com/WebChads/AccountService/internal/usecase"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"
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

func noOpMapper(s string) string { return s }

func NewDB(ctx context.Context, connString string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		return nil, err
	}

	db.MapperFunc(noOpMapper)

	return db, nil
}

func InitRouter(config *config.ServerConfig, logger *slog.Logger, db *sqlx.DB) http.Handler {
	rout := chi.NewRouter()
	http.Handle("/", rout)

	repos := usecase.NewRepositories(db)
	
	// Add all routers here
	accountUsecase := usecase.NewAccountUsecase(repos.Account, logger)
	accountRouter := router.NewAccountRouter(rout, config, logger, accountUsecase)
	// ...
	
	// Configure routers
	router.ConfigureAccountRouter(accountRouter)
	// ...
	
	// Serve Swagger UI
	rout.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	return rout
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
