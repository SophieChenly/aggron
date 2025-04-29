package main

import (
	"aggron/internal/api"
	"aggron/internal/config"
	"aggron/internal/db"
	"aggron/internal/repository"
	"aggron/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnvVariables()

	dbInstance, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	// init repositories (database interfaces)
	userRepo := repository.NewUserRepository(dbInstance)

	// init services
	s3Service, err := services.NewS3(services.S3Config{
		Region:          os.Getenv("AWS_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_ACCESS_KEY_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// init handlers
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 24 // 8 Mib

	router.POST("/file", api.UploadFile)
	router.GET("/file", api.RetrieveFile)
	router.POST("/auth/callback", api.CallbackAuth)

	router.Run(":8080")
}
