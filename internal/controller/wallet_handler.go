package controller

import (
	"net/http"
	"strconv"

	"github.com/yyq1025/balance-backend/internal/entity"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletService entity.WalletUseCase
}

func NewWalletHandler(r *gin.Engine, w entity.WalletUseCase) {
	handler := &WalletHandler{w}
	r.GET("/wallets", handler.GetManyWithPagination)
	r.GET("/wallets/:id", handler.GetOne)
	r.POST("/wallets", handler.AddOne)
	r.DELETE("/wallets/:id", handler.DeleteOne)
}

func (w *WalletHandler) GetManyWithPagination(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	idLte, err := strconv.Atoi(c.DefaultQuery("idLte", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid idLte"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "6"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid page size"})
		return
	}

	p := &entity.Pagination{IDLte: idLte, Page: page, PageSize: pageSize}

	balances, p, err := w.WalletService.GetManyWithPagination(c.Request.Context(), entity.Wallet{UserID: userID}, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balances": balances, "next": p})
}

func (w *WalletHandler) GetOne(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	balance, err := w.WalletService.GetOne(c.Request.Context(), entity.Wallet{UserID: userID, ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (w *WalletHandler) AddOne(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	wallet := entity.Wallet{UserID: userID}
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	wallet.NetworkName = wallet.Network.Name

	balance, err := w.WalletService.AddOne(c.Request.Context(), &wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (w *WalletHandler) DeleteOne(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	err = w.WalletService.DeleteOne(c.Request.Context(), entity.Wallet{UserID: userID, ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
