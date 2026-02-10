package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/middleware"
	"github.com/gin-gonic/gin"
)

func (a *App) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", a.Handlers.Auth.Login)
		auth.GET("/refresh", a.Handlers.Auth.RefreshToken)
		// Add more auth routes here as needed
		// auth.POST("/logout", a.Handlers.Auth.Logout)
	}
	user := r.Group("/user")
	{
		user.POST("/", a.Handlers.User.Register)
		user.GET("/:id", middleware.RequiredAuthentication(), a.Handlers.User.GetUser)
		// Add more user routes here as needed
		// user.GET("/:id", middleware.PasetoAuth(), a.Handlers.User.GetUser)
	}
}
