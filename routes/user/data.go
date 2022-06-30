package user

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// user_cache: email:User
var user_cache sync.Map

func CreateUser(db *gorm.DB, user *User) (int64, error) {
	result := db.Create(user)
	return result.RowsAffected, result.Error
}

func QueryUsers(db *gorm.DB, condition *User, users *[]User) (int64, error) {
	if cached_user, exist := user_cache.Load(condition.Email); exist {
		*users = []User{cached_user.(User)}
		return 1, nil
	}
	result := db.Where(condition).Find(users)
	for _, user := range *users {
		user_cache.LoadOrStore(user.Email, user)
	}
	return result.RowsAffected, result.Error
}

func UpdateUsers(db *gorm.DB, old *User, new *User) (int64, error) {
	result := db.Where(old).Updates(new)
	if result.RowsAffected > 0 {
		user_cache.Delete(old.Email)
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
