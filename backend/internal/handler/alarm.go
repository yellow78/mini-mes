package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/repository"
	"github.com/yellow78/mini-mes/backend/internal/service"
)

// AlarmHandler 告警 HTTP 處理器
type AlarmHandler struct {
	alarmRepo repository.AlarmRepository
	spcRepo   repository.SpcRepository
	equipSvc  *service.EquipmentService
}

func NewAlarmHandler(
	alarmRepo repository.AlarmRepository,
	spcRepo repository.SpcRepository,
	equipSvc *service.EquipmentService,
) *AlarmHandler {
	return &AlarmHandler{alarmRepo: alarmRepo, spcRepo: spcRepo, equipSvc: equipSvc}
}

// ListAlarms GET /api/v1/alarms
// ?all=true 回傳全部（含已確認），預設只回未確認
func (h *AlarmHandler) ListAlarms(c *gin.Context) {
	onlyUnacked := c.Query("all") != "true"
	alarms, err := h.alarmRepo.FindAll(c.Request.Context(), onlyUnacked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if alarms == nil {
		alarms = []model.AlarmEvent{}
	}
	c.JSON(http.StatusOK, gin.H{"data": alarms, "total": len(alarms)})
}

// AcknowledgeAlarm PUT /api/v1/alarms/:id/acknowledge
func (h *AlarmHandler) AcknowledgeAlarm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.alarmRepo.Acknowledge(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

// GetSpc GET /api/v1/spc/:equipment_id
func (h *AlarmHandler) GetSpc(c *gin.Context) {
	equipID, err := strconv.Atoi(c.Param("equipment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment_id"})
		return
	}
	limit := 100
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	records, err := h.spcRepo.FindByEquipment(c.Request.Context(), equipID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": records, "total": len(records)})
}
