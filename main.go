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

	if utils.GetEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	allowedOrigin := utils.GetEnv("ALLOWED_ORIGIN", "https://roast-myresume.vercel.app")
	log.Printf("Configuring CORS for origin: %s", allowedOrigin)

	corsConfig := cors.Config{
		AllowOrigins:     []string{allowedOrigin, "https://roast-myresume.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Serve static frontend files
	r.Static("/", "./dist")

	// React Router fallback
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Preflight request for /analyze
	r.OPTIONS("/analyze", func(c *gin.Context) {
		c.Status(200)
	})

	// Main API route
	r.POST("/analyze", handlers.AnalyzeResume)

	r.SetTrustedProxies(nil)

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	r.Run(":" + port)
}
