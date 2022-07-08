package user

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateUser(rc_cache *cache.Cache, db *gorm.DB, user *User) error {
	err := db.Create(user).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("user:%s", user.Email),
			Value: *user,
			TTL:   time.Hour,
		})
	}
	return err
}

func QueryUser(rc_cache *cache.Cache, db *gorm.DB, condition *User, user *User) error {
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("user:%s", condition.Email), user); err == nil {
		return nil
	}
	err := db.Where(condition).First(user).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("user:%s", user.Email),
			Value: *user,
			TTL:   time.Hour,
		})
	}
	return err
}

func UpdateUsers(rc_cache *cache.Cache, db *gorm.DB, old *User, new *User, users *[]User) error {
	err := db.Model(users).Clauses(clause.Returning{}).Where(old).Updates(new).Error
	for _, user := range *users {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("user:%s", user.Email),
			Value: user,
			TTL:   time.Hour,
		})
	}
	return err
}

func VerifyCode(rc_cache *cache.Cache, email string, code string) bool {
	var actual string
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("code:%s", email), &actual); err == nil && actual == code {
		rc_cache.Delete(context.TODO(), fmt.Sprintf("code:%s", email))
		return true
	}
	return false
}
