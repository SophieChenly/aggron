package main

import (
	"aggron/internal/config"
	"aggron/internal/db"
	"aggron/internal/repository"
	"aggron/internal/services"
	"log"
	"os"
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
}
