package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/middlewares"
	"github.com/saifwork/mock-service/internal/services"
)

type ProjectHandler struct {
	service *services.ProjectService
	cfg     *config.Config
}

func NewProjectHandler(service *services.ProjectService, cfg *config.Config) *ProjectHandler {
	return &ProjectHandler{service: service, cfg: cfg}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.RouterGroup) {
	projectRoutes := r.Group("/api/projects")
	projectRoutes.Use(middlewares.AuthMiddleware(h.cfg)) // JWT-based middleware
	{
		projectRoutes.POST("", h.CreateProject)
		projectRoutes.GET("/:pid", h.GetProjectByID)
		projectRoutes.GET("", h.GetProjects)
		projectRoutes.PUT("/:pid", h.UpdateProject)
		projectRoutes.DELETE("/:pid", h.DeleteProject)
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var payload struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	project, err := h.service.CreateProject(userID.(string), payload.Name, payload.Description)
	if err != nil {
		responses.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusCreated, "Project created successfully", project)
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	projectID := c.Param("pid")
	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	project, err := h.service.GetProjectByID(projectID, userID.(string))
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Project fetched successfully", project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projects, err := h.service.GetUserProjects(userID.(string))
	if err != nil {
		responses.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Projects fetched successfully", projects)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("pid")

	var payload struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.UpdateProject(projectID, userID.(string), payload.Name, payload.Description); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Project updated successfully", nil)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("pid")
	userID, exists := c.Get("userID")
	if !exists {
		responses.JSONError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.DeleteProject(projectID, userID.(string)); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Project deleted successfully", nil)
}
