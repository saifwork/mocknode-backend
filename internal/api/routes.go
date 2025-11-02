package api

import (
	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/handlers"
	"github.com/saifwork/mock-service/internal/core/config"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(
	r *gin.Engine,
	cfg *config.Config,

	authHandler *handlers.AuthHandler,
	projectHandler *handlers.ProjectHandler,
	collectionHandler *handlers.CollectionHandler,
	recordHandler *handlers.RecordHandler,
	healthHandler *handlers.HealthHandler,
) {
	// Handlers

	authHandler.RegisterRoutes(&r.RouterGroup)
	projectHandler.RegisterRoutes(&r.RouterGroup)
	collectionHandler.RegisterRoutes(&r.RouterGroup)
	recordHandler.RegisterRoutes(&r.RouterGroup)
	healthHandler.RegisterRoutes(&r.RouterGroup)
}
