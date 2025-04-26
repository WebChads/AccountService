package router

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/models/dtos"
	response "github.com/WebChads/AccountService/internal/pkg/api"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type AccountUsecase interface {
	Create(ctx context.Context, dto dtos.CreateAccountRequest) error
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

func ConfigureAccountRouter(r *AccountRouter) {
	// TODO: use middleware to check user authorization

	r.defaultHandler.HandleFunc("/api/v1/account/create-account", r.CreateAccountHandler)
	// ...
}

func (a *AccountRouter) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	var request dtos.CreateAccountRequest

	// Serialize account info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)

		// EOF means there is no data in the request body
		if errors.Is(err, io.EOF) {
			slog.Error("request body is empty")
			render.JSON(w, r, response.Error("empty request body"))
			return
		}

		slog.Error("failed to decode request body", slogerr.Error(err))
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}

	// Validate request fields
	err = validator.New().Struct(request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// Get phone number from auth token
	request.PhoneNumber = "892612345678"

	err = a.usecase.Create(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			render.Status(r, http.StatusGatewayTimeout)
			render.JSON(w, r, response.Error(err.Error()))
		}

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
}
