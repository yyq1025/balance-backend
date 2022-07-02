package wallet

import (
	"net/http"
	"strconv"

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

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

	address := data["address"]
	if address != "" && !utils.IsValidAddress(address) {
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

	Id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	res := DeleteWalletsByIds(db, &Wallet{Id: Id, UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	// data := make(map[string]string)

	// c.ShouldBindJSON(&data)

	// address := data["address"]
	// if address != "" && !utils.IsValidAddress(address) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "invalid address"})
	// 	return
	// }

	// network := data["network"]

	// token := data["token"]
	// if token != "" && !utils.IsValidAddress(token) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
	// 	return
	// }

	// tag := data["tag"]

	res := GetBalanceByParams(db, &Wallet{Id: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}
