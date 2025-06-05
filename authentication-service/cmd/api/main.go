package main

import (
	"auth/data"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type AuthConfig struct {
	Models data.Models
}

func main() {
	log.Println("Starting authentication service at port: ", webPort)

	app := AuthConfig{
		Models: data.New(),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
