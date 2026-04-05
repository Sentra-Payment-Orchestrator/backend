package routes

import (
	"net/http"

	"github.com/dwikie/sentra-payment-orchestrator/handler"
	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRoute struct {
	user_handler *handler.UserHandler
}

func InitUserRoute(pool *pgxpool.Pool) *UserRoute {
	return &UserRoute{user_handler: handler.NewUserHandler(pool, &handler.UserHandlerDependencies{})}
}

func (r *UserRoute) Createuser(c *gin.Context) {
	ctx := c.Request.Context()
	payload := model.CreateUserpayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.user_handler.InsertUser(ctx, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (r *UserRoute) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Implement logic to retrieve user by ID
	c.JSON(http.StatusOK, gin.H{"message": "GetUser endpoint", "user_id": userID})
}
