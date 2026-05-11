package main

import (
	"github.com/dwikie/sentra-payment-orchestrator/middleware"
	"github.com/dwikie/sentra-payment-orchestrator/routes"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Routes struct {
	Auth *routes.AuthRoute
	User *routes.UserRoute
}

func (a *App) RegisterRoutes(g *gin.Engine) {
	r := InitializeRoutes(a.Pool, a.Redis)

	auth_group := g.Group("/auth")
	{
		auth_group.POST("/login", r.Auth.Login)
		auth_group.GET("/refresh", r.Auth.RefreshToken)
	}

	user := routes.InitUserRoute(a.Pool)
	user_group := g.Group("/user")
	{
		user_group.POST("/", user.Createuser)
		user_group.GET("/", user.GetUser)
		user_group.GET("/:id", middleware.RequiredAuthentication(), user.GetUser)
		// Add more user routes here as needed
		// user.GET("/:id", middleware.PasetoAuth(), a.Handlers.User.GetUser)
	}
}

func InitializeRoutes(pool *pgxpool.Pool, redis *redis.Client) *Routes {
	auth := routes.InitAuthRoute(pool, redis)
	user := routes.InitUserRoute(pool)

	return &Routes{
		Auth: auth,
		User: user,
	}
}
