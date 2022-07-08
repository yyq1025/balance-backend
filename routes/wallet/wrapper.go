package wallet

import (
	"log"
	"net/http"
	"strconv"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
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

	res := AddWallet(rc_cache, db, &wallet)

	c.JSON(res.Code, res.Data)
	log.Print(rc_cache.Stats())
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
	log.Print(rc_cache.Stats())
}

func GetBalancesHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	res := GetBalances(rc_cache, db, &Wallet{UserId: userId})

	c.JSON(res.Code, res.Data)
	log.Print(rc_cache.Stats())
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

	res := GetBalance(rc_cache, db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
	log.Print(rc_cache.Stats())
}
