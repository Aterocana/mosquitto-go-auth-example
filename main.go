package main

import (
	"auth-server/server"
	"log"
)

func main() {
	srv := server.New()
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
