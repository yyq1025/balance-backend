package network

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetNetworksHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	res := GetAllNetWorks(db)

	c.JSON(res.Code, res.Data)
}
