package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/models"
	"github.com/saifwork/mock-service/internal/services"
)

type ConfigHandler struct {
	service *services.ConfigService
	cfg     *config.Config
}

func NewConfigHandler(service *services.ConfigService, cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{service: service, cfg: cfg}
}

func (h *ConfigHandler) RegisterRoutes(r *gin.RouterGroup) {
	configRoutes := r.Group("/config")
	{
		configRoutes.GET("/field-types", h.GetFieldTypes)

		// static data routes - for user, products, carts etc .
		configRoutes.GET("preset/:category", h.GetPresets)
		configRoutes.GET("preset/:category/:id", h.GetPresets)

	}
}

func (h *ConfigHandler) GetFieldTypes(c *gin.Context) {
	fieldTypes := h.service.GetFieldTypes()
	responses.JSONSuccess(c, http.StatusOK, "Field types fetched successfully", gin.H{"fieldTypes": fieldTypes})
}

func (h *ConfigHandler) GetPresets(c *gin.Context) {
	category := c.Param("category")
	id := c.Param("id")

	var data []models.PresetRecord
	switch category {
	case "users":
		data = models.GetPresetUsers()
	case "products":
		data = models.GetPresetProducts()
	case "comments":
		data = models.GetPresetComments()
	case "carts":
		data = models.GetPresetCarts()
	case "posts":
		data = models.GetPresetPosts()
	case "todos":
		data = models.GetPresetTodos()
	default:
		responses.JSONError(c, http.StatusNotFound, "Invalid preset category")
		return
	}

	if id == "" {
		responses.JSONSuccess(c, http.StatusOK, "Preset records fetched", data)
		return
	}

	for _, record := range data {
		if record["id"] == id {
			responses.JSONSuccess(c, http.StatusOK, "Preset record found", record)
			return
		}
	}

	responses.JSONError(c, http.StatusNotFound, "Record not found")
}
