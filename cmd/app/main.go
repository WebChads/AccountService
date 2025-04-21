package main

import (
	server "github.com/WebChads/AccountService/internal/delivery/http"
	storage "github.com/WebChads/AccountService/internal/storage/pgsql/account"
)

// TODO: init server configuration
// TODO: setup logger
// TODO: init database
// TODO: pkg: auth middleware, logger
// TODO: swagger docs

func main() {
	// init config
	config, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// init database
	db := server.NewDB(config)

	// init repos, usecases and API handlers
	repos := storage.NewAccountRepository()

	// setup logger

	// run server
	srv := server.NewServer(8080)

	srv.ListenAndServe()
}
