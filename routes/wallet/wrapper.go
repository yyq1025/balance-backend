package wallet

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func CreateWalletHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	wallet := Wallet{UserID: userID}
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	wallet.NetworkName = wallet.Network.Name

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := addWallet(ctx, rdbCache, db, &wallet)

	c.JSON(res.Code, res.Data)
}

func DeleteWalletHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

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

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := DeleteBalance(ctx, rdbCache, db, &Wallet{ID: id, UserID: userID})

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

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

	p := Pagination{IDLte: idLte, Page: page, PageSize: pageSize}

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := getBalancesWithPagination(ctx, rdbCache, db, &Wallet{UserID: userID}, &p)

	c.JSON(res.Code, res.Data)
}

func GetBalanceHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userID := c.MustGet("userID").(string)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := getBalance(ctx, rdbCache, db, &Wallet{ID: id, UserID: userID})

	c.JSON(res.Code, res.Data)
}
