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
	rdb_cache := c.MustGet("rdb_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(string)
	if userId == "" {
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

	tag := data["tag"]

	wallet := Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: network,
		Token:   common.HexToAddress(token),
		Tag:     tag,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := AddWallet(ctx, rdb_cache, db, &wallet)

	c.JSON(res.Code, res.Data)
}

func DeleteWalletsHandler(c *gin.Context) {
	rdb_cache := c.MustGet("rdb_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(string)
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	res := DeleteBalances(rdb_cache, db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	rdb_cache := c.MustGet("rdb_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(string)
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := GetBalances(ctx, rdb_cache, db, &Wallet{UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalanceHandler(c *gin.Context) {
	rdb_cache := c.MustGet("rdb_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(string)
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := GetBalance(ctx, rdb_cache, db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}
