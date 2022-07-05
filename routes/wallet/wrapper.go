package wallet

import (
	"net/http"
	"strconv"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/common"
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

	wallet := Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: network,
		Token:   common.HexToAddress(token),
		Tag:     tag,
	}

	res := AddWallet(db, &wallet)

	c.JSON(res.Code, res.Data)
}

func DeleteWalletsHandler(c *gin.Context) {
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

	res := DeleteBalances(db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalancesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	res := GetBalances(db, &Wallet{UserId: userId})

	c.JSON(res.Code, res.Data)
}

func GetBalanceHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userId := c.MustGet("userId").(int)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	res := GetBalance(db, &Wallet{ID: id, UserId: userId})

	c.JSON(res.Code, res.Data)
}
