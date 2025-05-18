package main

import (
	"log"
	"roastmyresume/handlers"
	"roastmyresume/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()

	r := gin.Default()

	allowedOrigin := utils.GetEnv("ALLOWED_ORIGIN", "http://localhost:5173")
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
	r.Run(":" + utils.GetEnv("PORT", "8080"))
}
