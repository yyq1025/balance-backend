package controller

import (
	"net/http"
	"yyq1025/balance-backend/internal/entity"

	"github.com/gin-gonic/gin"
)

type NetworkHandler struct {
	NetworkService entity.NetworkUseCase
}

func NewNetworkHandler(r *gin.Engine, n entity.NetworkUseCase) {
	handler := &NetworkHandler{n}
	r.GET("/networks", handler.GetAll)
}

func (h *NetworkHandler) GetAll(c *gin.Context) {
	networks, err := h.NetworkService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"networks": networks})
}
