package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// create router group for handlers
func Register(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/status", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
}
