package utils

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	// dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=require", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("SQL Connected!")
	return db
}

func GetRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":6379",
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Redis Connected!")
	return client
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	return passwordvalidator.Validate(password, 60) == nil
}

func IsValidCode(code string) bool {
	return regexp.MustCompile("^[0-9]{6}$").MatchString(code)
}

func CreateJWT(userId int, duration time.Duration) (string, error) {
	claims := &Claims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(SecretKey))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(SecretKey), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
