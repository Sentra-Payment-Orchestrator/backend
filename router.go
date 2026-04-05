package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/middleware"
	"github.com/dwikie/sentra-payment-orchestrator/routes"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	auth := routes.InitAuthRoute(pool)
	auth_group := r.Group("/auth")
	{
		auth_group.POST("/login", auth.Login)
		auth_group.GET("/refresh", auth.RefreshToken)
	}

	user := routes.InitUserRoute(pool)
	user_group := r.Group("/user")
	{
		user_group.POST("/", user.Createuser)
		user_group.GET("/", user.GetUser)
		user_group.GET("/:id", middleware.RequiredAuthentication(), user.GetUser)
		// Add more user routes here as needed
		// user.GET("/:id", middleware.PasetoAuth(), a.Handlers.User.GetUser)
	}
}
