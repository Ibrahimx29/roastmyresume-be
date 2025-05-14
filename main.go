package main

import (
	"log"
	"os"
	"roastmyresume/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	log.Printf("Configuring CORS for origin: %s", allowedOrigin)

	corsConfig := cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

    r.POST("/analyze", handlers.AnalyzeResume)
    r.Run(":8080") // starts server at http://localhost:8080
}