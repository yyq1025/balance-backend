package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetNetworksHandler(c *gin.Context) {
	rc_cache := c.MustGet("rc_cache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	res := GetNetWorks(rc_cache, db, &Network{})

	c.JSON(res.Code, res.Data)
}
