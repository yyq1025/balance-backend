package wallet

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/common"
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

	data := make(map[string]string)

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	address := data["address"]
	if !utils.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
		return
	}

	network := data["network"]

	token := data["token"]
	if token != "" && !utils.IsValidAddress(token) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}

	wallet := Wallet{
		UserID:  userID,
		Address: common.HexToAddress(address),
		Network: network,
		Token:   common.HexToAddress(token),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := AddWallet(ctx, rdbCache, db, &wallet)

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

	res := GetBalancesWithPagination(ctx, rdbCache, db, &Wallet{UserID: userID}, &p)

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

	res := GetBalance(ctx, rdbCache, db, &Wallet{ID: id, UserID: userID})

	c.JSON(res.Code, res.Data)
}
