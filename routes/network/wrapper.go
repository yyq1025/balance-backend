package network

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetNetworksHandler(c *gin.Context) {
	rdbCache := c.MustGet("rdbCache").(*cache.Cache)

	db := c.MustGet("db").(*gorm.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()

	res := getNetWorks(ctx, rdbCache, db, &Network{})

	c.JSON(res.Code, res.Data)
}
