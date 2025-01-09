package main

import (
	"bima-go/config"
	"bima-go/models"
	"bima-go/routes"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	appPort := os.Getenv("APP_PORT")
	models.ConnectDatabase()
	router := routes.SetupRouter()
	err := router.Run(":" + appPort)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
