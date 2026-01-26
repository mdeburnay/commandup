package main

import (
	"log"

	"commandup/config"
	"commandup/models"
	"commandup/routers"
	"net/http"
)

func init() {
	models.Init()
	config.Init()
}

func main() {
	routersInit := routers.InitRouter()
	port := ":8080"

	server := &http.Server{
		Addr:    port,
		Handler: routersInit,
	}

	log.Printf("Server started at %s", port)
	server.ListenAndServe()
}
