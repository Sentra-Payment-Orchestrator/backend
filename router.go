package main

import "github.com/gin-gonic/gin"

func (a *App) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("auth")
	{
		auth.POST("/login", a.Handlers.Auth.Login)
		auth.POST("/register", a.Handlers.Auth.Register)
		// Add more auth routes here as needed
		// auth.POST("/logout", a.Handlers.Auth.Logout)
		// auth.POST("/register", a.Handlers.Auth.Register)
	}
}
