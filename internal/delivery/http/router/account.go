package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/WebChads/AccountService/internal/config"
	"github.com/WebChads/AccountService/internal/models/dtos"
	response "github.com/WebChads/AccountService/internal/pkg/api"
	slogerr "github.com/WebChads/AccountService/internal/pkg/logger"
	"github.com/WebChads/AccountService/pkg/middleware/auth"
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
	jwt := auth.JWTConfig{
		SecretKey: r.config.SecretKey,
	}

	// Auth middleware
	r.defaultHandler.Use(jwt.AuthMiddleware)

	r.defaultHandler.HandleFunc("/api/v1/account/create-account", r.CreateAccountHandler)
	// ...
}

// CreateAccountHandler godoc
// @Title CreateAccount
// @Summary Create a new personal account
// @Description This endpoint creates a new user personal account
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dtos.CreateAccountRequest true "Create account parameters"
// @Success 200 {object} dtos.Response "Successfully created personal account"
// @Failure 400 {object} dtos.Response "Request body is empty or Request body field validation failed"
// @Router /api/v1/account/create-account [post]
func (a *AccountRouter) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	var request dtos.CreateAccountRequest

	// Serialize account info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		// render.Status(r, http.StatusBadRequest)

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

	// Get phone number from request context
	// request.PhoneNumber = r.Context().Value("phone_number").(string)

	// Use mock data for now
	request.PhoneNumber = "892612345678"

	err = a.usecase.Create(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			response.JSON(w, http.StatusRequestTimeout, err.Error())
		}

		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
}

func getValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	}

	return "validation error"
}
