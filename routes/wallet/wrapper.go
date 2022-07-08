package wallet

import (
	"net/http"
	"strconv"
	"time"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func CreateWalletHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

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

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	res := AddWallet(ctx, rc_cache, db, &wallet)

	c.JSON(res.Code, res.Data)
}

func DeleteWalletsHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	res := DeleteBalances(rc_cache, db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	res := GetBalances(ctx, rc_cache, db, &Wallet{UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalanceHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	res := GetBalance(ctx, rc_cache, db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}
