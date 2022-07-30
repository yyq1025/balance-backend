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
	// ctx, cancel := context.WithTimeout(c.Request.Context(), 3500*time.Millisecond)
	// defer cancel()
	networks, err := h.NetworkService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(getStatusCode(err), ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"networks": networks})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case entity.ErrGetNetwork:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
