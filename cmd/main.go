package main

import (
	"auth/domain"
	"auth/pkg/utils"
	"auth/presentation/controller"
	"auth/repository"
	"auth/service"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
	databaseURL := os.Getenv("DATABASE_URL")

	if jwtSecret == "" || clientID == "" || clientSecret == "" || clientCallbackURL == "" || databaseURL == "" {
		err := fmt.Sprintf("Missing required environment variables: JWT_SECRET=%s, CLIENT_ID=%s, CLIENT_SECRET=%s, CLIENT_CALLBACK_URL=%s, DATABASE_URL=%s", jwtSecret, clientID, clientSecret, clientCallbackURL, databaseURL)
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&domain.User{})

	jwtUtils := utils.NewJWTUtil(jwtSecret)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtUtils)
	authController := controller.NewAuthController(userService)
	oauthController := controller.NewOAuthController(userService)

	goth.UseProviders(
		google.New(
			clientID,
			clientSecret,
			clientCallbackURL,
		),
	)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	oauth := r.Group("/oauth")
	{
		oauth.GET("/auth/:provider", oauthController.BeginAuth)
		oauth.GET("/callback/:provider", oauthController.CallbackHandler)
	}

	//TODO read port from env
	r.Run(":5000")
}
