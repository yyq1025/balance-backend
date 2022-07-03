package user

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// user_cache: email:User
var user_cache sync.Map

func CreateUser(db *gorm.DB, user *User) error {
	err := db.Create(user).Error
	if err == nil {
		user_cache.Store(user.Email, *user)
	}
	return err
}

func QueryUser(db *gorm.DB, condition *User, user *User) error {
	if cached_user, exist := user_cache.Load(condition.Email); exist {
		*user = cached_user.(User)
		return nil
	}
	err := db.Where(condition).First(user).Error
	if err == nil {
		user_cache.Store(user.Email, *user)
	}
	return err
}

func UpdateUsers(db *gorm.DB, old *User, new *User, users *[]User) (int64, error) {
	result := db.Model(users).Clauses(clause.Returning{}).Where(old).Updates(new)
	for _, user := range *users {
		user_cache.Store(user.Email, user)
	}
	return result.RowsAffected, result.Error
}

func VerifyCode(rc *redis.Client, email string, code string) bool {
	actual, err := rc.Get(context.Background(), email).Result()
	if err == nil && code == actual {
		rc.Del(context.Background(), email)
		return true
	}
	return false
}
