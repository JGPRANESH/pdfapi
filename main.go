package main

import (
	"log"
	"os"

	"pdfapi/config"
	"pdfapi/routes"
	"pdfapi/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	config.InitFirebase()

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		panic(err)
	}

	services.StartCleanupWorker()

	router := gin.Default()

	router.MaxMultipartMemory = 10 << 20

	routes.SetupRoutes(router)

	router.Run(":8080")
}
