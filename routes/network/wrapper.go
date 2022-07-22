package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetNetworksHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	res := GetNetWorks(rdbCache, db, &Network{})

	c.JSON(res.Code, res.Data)
}
