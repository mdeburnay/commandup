package main

import (
	"log"

	"commandup/models"
	"commandup/routers"
	"net/http"
)

func init() {
	models.Init()
}

func main() {
	routersInit := routers.InitRouter()
	port := "localhost:8080"

	server := &http.Server{
		Addr:    port,
		Handler: routersInit,
	}

	log.Printf("Server started at %s", port)
	server.ListenAndServe()
}
