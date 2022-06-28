package network

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NetworksHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	res := GetAllNetWorks(db)

	c.JSON(res.Code, res.Data)
}

func NetworkByNameHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	network := c.Param("network")

	res := GetNetworkInfoByName(db, network)

	c.JSON(res.Code, res.Data)
}
