package main

import (
	"aggron/internal/api"
	"aggron/internal/cache"
	"aggron/internal/config"
	"aggron/internal/db"
	"aggron/internal/repository"
	"aggron/internal/services"
	"log"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {

	config.LoadEnvVariables()

	dbInstance, err := db.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	// init repositories (database interfaces)
	// users
	userRepo := repository.NewUserRepository(dbInstance)

	// init services
	// s3
	s3Service, err := services.NewS3(services.S3Config{
		Region:          os.Getenv("AWS_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_ACCESS_KEY_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// auth
	authService, err := services.NewAuth(services.AuthConfig{
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OIDC_REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "email"},
		IssuerURL:    os.Getenv("OIDC_ISSUER_URL"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// redis (cache)
	redisService := cache.NewRedis(*redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 2, // connection
	}))

	// init handlers
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 24 // 8 Mib

	fileHandler := api.FileController{AuthService: authService, RedisService: redisService}
	authHandler := api.AuthController{AuthService: authService, RedisService: redisService}

	router.POST("/file", fileHandler.UploadFile)
	router.GET("/file", fileHandler.RetrieveFile)
	router.GET("/auth/callback", authHandler.AuthCallback)

	router.Run(":8080")
}
