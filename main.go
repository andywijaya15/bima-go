package main

import (
	"bima-go/config"
	"bima-go/models"
	"bima-go/routes"
	"log"

	"github.com/gin-contrib/gzip"
)

func main() {

	config.LoadEnv()
	router := routes.SetupRouter()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	models.ConnectDatabase()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
