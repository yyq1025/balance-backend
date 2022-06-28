package wallet

import (
	"net/http"

	"yyq1025/balance-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateWalletHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	address := c.Query("address")
	if !utils.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
		return
	}

	network := c.Query("network")

	token := c.Query("token")
	if token != "" && !utils.IsValidAddress(token) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}

	tag := c.Query("tag")

	res := AddWallet(db, userId, address, network, token, tag)

	c.JSON(res.Code, res.Data)
}

func GetWalletsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	address := c.Query("address")
	if address != "" && !utils.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
		return
	}

	network := c.Query("network")

	token := c.Query("token")
	if token != "" && !utils.IsValidAddress(token) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}

	tag := c.Query("tag")

	res := GetWalletsByParams(db, userId, address, network, token, tag)

	c.JSON(res.Code, res.Data)
}

func DeleteWalletsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	address := c.Query("address")
	if address != "" && !utils.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
		return
	}

	network := c.Query("network")

	token := c.Query("token")
	if token != "" && !utils.IsValidAddress(token) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}

	tag := c.Query("tag")

	res := DeleteWalletsByParams(db, userId, address, network, token, tag)

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	address := c.Query("address")
	if address != "" && !utils.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
		return
	}

	network := c.Query("network")

	token := c.Query("token")
	if token != "" && !utils.IsValidAddress(token) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}

	tag := c.Query("tag")

	res := GetBalanceByParams(db, userId, address, network, token, tag)

	c.JSON(res.Code, res.Data)
}
