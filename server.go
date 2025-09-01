package main

import (
	"backend/handlers"
	"backend/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	authHandler := handlers.NewAuthHandler(DB, jwtSecret)
	guideHandler := handlers.NewGuideHandler(DB)
	bookingHandler := handlers.NewBookingHandler(DB)

	// Public
	router.GET("/health", func(context *gin.Context) {
		context.JSON(200, gin.H{"ok": true})
	})
	router.POST("/auth/register", authHandler.RegisterUser)
	router.POST("/auth/login", authHandler.LoginUser)
	router.GET("/guides", guideHandler.ListGuides)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.GET("/bookings", bookingHandler.ListUserBookings)
		auth.POST("/bookings", bookingHandler.CreateBooking)
	}

	return router
}
