package main

import (
	"net/http"

	"github.com/jipark0716/dalmuti/routers"
)

func main() {
	server := &http.Server{
		Handler: routers.InitRouter(),
	}

	server.ListenAndServe()
}
