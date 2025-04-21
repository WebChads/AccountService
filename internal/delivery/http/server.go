package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/delivery/http/router"
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

func NewDB(config *config.ServerConfig) *sql.DB {
	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		return nil
	}

	return db	
}

func InitRouter(config *config.ServerConfig, logger *slog.Logger) *chi.Mux {
	// router := chi.NewRouter()
	// accountRepo := storage.NewAccountRepository()
	// accountUsecase := usecase.NewAccountUsecase(accountRepo)
	// account := handler.NewAccountHandler(usecase.NewAccountUsecase())
	// router.HandleFunc("/api/v1/account/create-account", account.CreateAccountHandler)

	rout := chi.NewRouter()

	// add all routers here
	accountRouter := router.NewAccountRouter(rout, config, logger)

	// Routers
	// pingrouter := PingRouter.NewPingRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), UserUC, log)
	// instructionrouter := InstructionRouter.NewInstructionRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), InstructionUC, UserUC, log)
	// productrouter := ProductRouter.NewProductRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), ProductUC, UserUC, log)
	// rentedproductrouter := RentedProductRouter.NewRentedProductRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), RentedProductUC, UserUC, log)
	// showcaserouter := ShowcaseRouter.NewShowcaseRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), ShowcaseUC, UserUC, log)
	// userrouter := UserRouter.NewUserRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), rout.PathPrefix("/api/v1").Subrouter(), UserUC, ObjectUC, log)
	// objectrouter := ObjectRouter.NewObjectRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), ObjectUC, UserUC, log)
	// cardrouter := CardRouter.NewCardRouter(s.config, rout.PathPrefix("/api/v1").Subrouter(), CardUC, UserUC, log)

	http.Handle("/", rout)

	// Configure Routers
	router.ConfigureAccountRouter(accountRouter)
	// PingRouter.ConfigureRouter(pingrouter)
	// InstructionRouter.ConfigureRouter(instructionrouter)
	// ProductRouter.ConfigureRouter(productrouter)
	// RentedProductRouter.ConfigureRouter(rentedproductrouter)
	// ShowcaseRouter.ConfigureRouter(showcaserouter)
	// UserRouter.ConfigureRouter(userrouter)
	// ObjectRouter.ConfigureRouter(objectrouter)
	// CardRouter.ConfigureRouter(cardrouter)

	return router
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
