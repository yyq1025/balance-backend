package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetNetworksHandler(c *gin.Context) {
	rdb_cache := c.MustGet("rdb_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	res := GetNetWorks(rdb_cache, db, &Network{})

	c.JSON(res.Code, res.Data)
}
