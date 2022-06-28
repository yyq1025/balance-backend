package user

import (
	"net/http"

	"yyq1025/balance-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func RegisterHandler(c *gin.Context) {
	rc := c.MustGet("rc").(*redis.Client)

	db := c.MustGet("db").(*gorm.DB)

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

	email := data["email"]
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email"})
		return
	}

	password := data["password"]
	if !utils.IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	code := data["code"]
	if !utils.IsValidCode(code) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid code"})
		return
	}

	res := AddUser(rc, db, email, password, code)

	c.JSON(res.Code, res.Data)
}

func SendCodeHandler(c *gin.Context) {
	rc := c.MustGet("rc").(*redis.Client)

	sender := c.MustGet("sender").(*utils.Sender)

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

	email := data["email"]
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email"})
		return
	}

	res := SendCode(sender, rc, email)

	c.JSON(res.Code, res.Data)
}

func LoginHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

	email := data["email"]
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email"})
		return
	}

	password := data["password"]
	if !utils.IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	res := Login(db, email, password)

	c.JSON(res.Code, res.Data)
}

func ChangePasswordHandler(c *gin.Context) {
	rc := c.MustGet("rc").(*redis.Client)

	db := c.MustGet("db").(*gorm.DB)

	data := make(map[string]string)

	c.ShouldBindJSON(&data)

	email := data["email"]
	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email"})
		return
	}

	password := data["password"]
	if !utils.IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	code := data["code"]
	if !utils.IsValidCode(code) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid code"})
		return
	}

	res := ChangePassword(rc, db, email, password, code)

	c.JSON(res.Code, res.Data)
}
