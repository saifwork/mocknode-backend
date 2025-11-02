package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saifwork/mock-service/internal/api/responses"
	"github.com/saifwork/mock-service/internal/core/config"
	"github.com/saifwork/mock-service/internal/middlewares"
	"github.com/saifwork/mock-service/internal/services"
)

type RecordHandler struct {
	service *services.RecordService
	cfg     *config.Config
}

func NewRecordHandler(service *services.RecordService, cfg *config.Config) *RecordHandler {
	return &RecordHandler{service: service, cfg: cfg}
}

func (h *RecordHandler) RegisterRoutes(r *gin.RouterGroup) {
	recordRoutes := r.Group("api/collections/:collectionId/records")
	recordRoutes.Use(middlewares.AuthMiddleware(h.cfg))
	{
		recordRoutes.POST("", h.CreateRecord)
		recordRoutes.GET("", h.GetRecordsByCollection)
		recordRoutes.GET("/:rid", h.GetRecordByID)
		recordRoutes.PUT("/:rid", h.UpdateRecord)
		recordRoutes.DELETE("/:rid", h.DeleteRecord)
	}
}

func (h *RecordHandler) CreateRecord(c *gin.Context) {
	collectionID := c.Param("collectionId")

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	record, err := h.service.CreateRecord(collectionID, data)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusCreated, "Record created", record)
}

func (h *RecordHandler) GetRecordsByCollection(c *gin.Context) {
	collectionID := c.Param("collectionId")

	records, err := h.service.GetRecordsByCollection(collectionID)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Records fetched", records)
}

func (h *RecordHandler) GetRecordByID(c *gin.Context) {
	rid := c.Param("rid")

	record, err := h.service.GetRecordByID(rid)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Record fetched", record)
}

func (h *RecordHandler) UpdateRecord(c *gin.Context) {
	rid := c.Param("rid")

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		responses.JSONError(c, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	record, err := h.service.UpdateRecord(rid, data)
	if err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Record updated", record)
}

func (h *RecordHandler) DeleteRecord(c *gin.Context) {
	rid := c.Param("rid")

	if err := h.service.DeleteRecord(rid); err != nil {
		responses.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	responses.JSONSuccess(c, http.StatusOK, "Record deleted", nil)
}
