package main

import (
	"bima-go/config"
	"bima-go/models"
	"bima-go/routes"
	"log"
)

func main() {
	config.LoadEnv()
	models.ConnectDatabase()
	router := routes.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
