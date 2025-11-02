package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/middlewares"
	"github.com/saifwork/mock-service/internal/models"
	"github.com/saifwork/mock-service/internal/services"
)

type CollectionHandler struct {
	service *services.CollectionService
	cfg     *config.Config
}

func NewCollectionHandler(service *services.CollectionService, cfg *config.Config) *CollectionHandler {
	return &CollectionHandler{service: service, cfg: cfg}
}

func (h *CollectionHandler) RegisterRoutes(r *gin.RouterGroup) {
	collectionRoutes := r.Group("/api/projects/:pid/collections")
	collectionRoutes.Use(middlewares.AuthMiddleware(h.cfg))
	{
		collectionRoutes.POST("", h.CreateCollection)
		collectionRoutes.GET("", h.GetCollectionsByProject)
		collectionRoutes.GET("/:cid", h.GetCollectionByID)
		collectionRoutes.DELETE("/:cid", h.DeleteCollection)
	}
}

func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	projectID := c.Param("projectId")

	var body struct {
		Name   string                   `json:"name" binding:"required"`
		Fields []models.FieldDefinition `json:"fields"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	collection, err := h.service.CreateCollection(projectID, body.Name, body.Fields)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusCreated, "Collection created", collection)
}

func (h *CollectionHandler) GetCollectionsByProject(c *gin.Context) {
	projectID := c.Param("projectId")

	collections, err := h.service.GetCollectionsByProject(projectID)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Collections fetched", collections)
}

func (h *CollectionHandler) GetCollectionByID(c *gin.Context) {
	cid := c.Param("cid")

	collection, err := h.service.GetCollectionByID(cid)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Collection fetched", collection)
}

func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	cid := c.Param("cid")

	if err := h.service.DeleteCollection(cid); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Collection deleted", nil)
}
