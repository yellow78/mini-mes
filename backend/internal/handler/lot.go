package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/service"
)

// LotHandler Lot HTTP 處理器
type LotHandler struct {
	lotSvc      *service.LotService
	dispatchSvc *service.DispatchService
	hub         Broadcaster
}

func NewLotHandler(lotSvc *service.LotService, dispatchSvc *service.DispatchService, hub Broadcaster) *LotHandler {
	return &LotHandler{lotSvc: lotSvc, dispatchSvc: dispatchSvc, hub: hub}
}

// ListLots GET /api/v1/lots
func (h *LotHandler) ListLots(c *gin.Context) {
	lots, err := h.lotSvc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lots, "total": len(lots)})
}

// GetLot GET /api/v1/lots/:id
func (h *LotHandler) GetLot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	lot, err := h.lotSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if lot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lot})
}

// CreateLot POST /api/v1/lots
func (h *LotHandler) CreateLot(c *gin.Context) {
	var lot model.Lot
	if err := c.ShouldBindJSON(&lot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.lotSvc.Create(c.Request.Context(), &lot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": lot})
}

// DispatchLot POST /api/v1/lots/:id/dispatch
func (h *LotHandler) DispatchLot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	result, err := h.dispatchSvc.Dispatch(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.hub.Broadcast("lot_dispatched", gin.H{
		"lot_id":       result.LotID,
		"equipment_id": result.EquipmentID,
	})
	// 設備狀態改為 RUNNING，通知 Dashboard 即時更新
	h.hub.Broadcast("equipment_status_changed", gin.H{
		"equipment_id": result.EquipmentID,
		"status":       "RUNNING",
	})
	c.JSON(http.StatusOK, gin.H{"data": result})
}
