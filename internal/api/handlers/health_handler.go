package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	db *mongo.Client
}

func NewHealthHandler(db *mongo.Client) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/health", h.CheckHealth)
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := h.db.Ping(ctx, nil); err != nil {
		dbStatus = "down"
	}

	responses.JSONSuccess(c, http.StatusOK, "Service health status", gin.H{
		"service": "up",
		"db":      dbStatus,
		"time":    time.Now().Format(time.RFC3339),
	})
}
