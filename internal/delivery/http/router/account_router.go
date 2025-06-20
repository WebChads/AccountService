package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/models/dtos"
	response "github.com/WebChads/AccountService/internal/pkg/api"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/WebChads/AccountService/pkg/middleware/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type AccountUsecase interface {
	Create(ctx context.Context, dto dtos.CreateAccountRequest) error
	Get(ctx context.Context, userId string) (*dtos.GetAccountResponse, error)
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
	// Auth middleware
	authMiddleware := auth.NewMiddleware(r.config.AuthServiceUrl)

	r.defaultHandler.With(authMiddleware.Handler).Post("/api/v1/account/create-account", r.CreateAccountHandler)
	r.defaultHandler.With(authMiddleware.Handler).Get("/api/v1/account/get-account/{user_id}", r.GetAccountHandler)
	// r.defaultHandler.Patch("/api/v1/account/update-account", r.UpdateAccountHandler)
	// ...
}

// @Title GetAccount
// @Summary Get user account by ID
// @Description Returns account information for specified user ID
// @Tags Account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path string true "User ID" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} dtos.GetAccountResponse
// @Failure 400 {object} dtos.Response
// @Failure 404 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /api/v1/account/get-account/{user_id} [get]
func (a *AccountRouter) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		a.logger.Error("user_id param is empty")

		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	account, err := a.usecase.Get(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, account)
}

// @Title CreateAccount
// @Summary Create new account
// @Description Creates a new user account with provided details
// @Tags Account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dtos.CreateAccountRequest true "Account creation data"
// @Success 201 {object} dtos.Response
// @Failure 400 {object} dtos.Response
// @Failure 500 {object} dtos.Response
// @Router /api/v1/account/create-account [post]
func (a *AccountRouter) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	var request dtos.CreateAccountRequest

	// Serialize account info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		// EOF means there is no data in the request body
		if errors.Is(err, io.EOF) {
			a.logger.Error("request body is empty", slogerr.Error(err))
			response.JSON(w, http.StatusBadRequest, "request body is empty")
			return
		}

		a.logger.Error("failed to decode request body", slogerr.Error(err))
		response.JSON(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	// Validate request fields
	err = validator.New().Struct(request)
	if err != nil {
		var errors []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errors = append(errors, getValidationMsg(fieldErr))
			}
		}

		resp := map[string]any{"errors": errors}

		response.JSON(w, http.StatusBadRequest, resp)
		return
	}

	// Get user id from request context
	userIdAsString := r.Context().Value("user_id").(string)
	request.UserId, err = uuid.Parse(userIdAsString)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "unable to parse uuid from user_id")
	}

	err = a.usecase.Create(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			response.JSON(w, http.StatusRequestTimeout, err.Error())
		} else if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}

func getValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	}

	return "validation error"
}
