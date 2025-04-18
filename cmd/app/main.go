package main

import server "github.com/WebChads/AccountService/internal/delivery/http"

func main() {
	// init config

	// setup logger

	// run server
	srv := server.NewServer(5000)

	srv.ListenAndServe()
}