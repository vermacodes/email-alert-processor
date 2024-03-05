package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vermacodes/email-alert-processor/internal/entity"
)

type EmailAlertProcessingHandler struct {
	alertService entity.AlertService
}

func NewEmailAlertProcessingHandler(r *gin.RouterGroup, alertService entity.AlertService) {
	handler := &EmailAlertProcessingHandler{
		alertService: alertService,
	}

	r.POST("/cnflnt/new-case-alert", handler.processAlert)
}

func (h *EmailAlertProcessingHandler) processAlert(c *gin.Context) {
	alert := entity.Alert{}
	if err := c.BindJSON(&alert); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	alert, err := h.alertService.ProcessAlert(alert)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, alert)
}
