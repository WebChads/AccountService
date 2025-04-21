package router

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/models/dtos"
	response "github.com/WebChads/AccountService/internal/pkg/api"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type AccountUsecase interface {
	Create(dtos.CreateAccountRequest) error
}

type AccountRouter struct {
	defaultHandler *chi.Mux
	logger         *slog.Logger
	config         *config.ServerConfig
	usecase        AccountUsecase
}

func NewAccountRouter(r *chi.Mux, cfg *config.ServerConfig,
	log *slog.Logger, usecase AccountUsecase) *AccountRouter {
	router := &AccountRouter{
		defaultHandler: r,
		logger:         log,
		config:         cfg,
		usecase:        usecase,
	}

	return router
}

func ConfigureAccountRouter(ar *AccountRouter) {
	// TODO: use middleware to check user authorization

	ar.defaultHandler.HandleFunc("/api/v1/account/create-account", ar.CreateAccountHandler)
}

func (a *AccountRouter) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateAccountHandler")

	var request dtos.CreateAccountRequest

	// serialize account info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		// EOF means there is no data in the request body
		if errors.Is(err, io.EOF) {
			slog.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request body"))
			return
		}

		slog.Error("failed to decode request body", slogerr.Error(err))

		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	err = a.usecase.Create(request)
	if err != nil {
		return
	}
}
