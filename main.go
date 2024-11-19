package main

import (
	"auth/entities"
	"auth/infrastructure"
	"auth/presentation"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" || jwtSecret == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL, JWT_SECRET) are required")
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)

	db, err := infrastructure.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", presentation.Home)
	r.GET("/auth/:provider", presentation.SignInWithProvider)
	r.GET("/auth/:provider/callback", presentation.CallbackHandler)
	r.GET("/success", presentation.Success)

	r.Run(":5000")
}
