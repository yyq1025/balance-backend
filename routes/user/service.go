package user

import (
	"log"
	"net/http"
	"time"

	"yyq1025/balance-backend/utils"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var sender = newSender()

func AddUser(rc *redis.Client, db *gorm.DB, email, password, code string) utils.Response {
	if !VerifyCode(rc, email, code) {
		return utils.VerificationCodeError
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return utils.CreateUserError
	}
	user := User{Email: email, Password: hashedPassword}
	if err := CreateUser(db, &user); err != nil {
		log.Print(err)
		return utils.CreateUserError
	}
	jwt, err := utils.CreateJWT(user.ID, 24*time.Hour)
	if err != nil {
		log.Print(err)
		return utils.UserLoginError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"email": user.Email, "token": jwt}}
}

func SendCode(rc *redis.Client, email string) utils.Response {
	if err := sender.sendCode(rc, email); err != nil {
		log.Print(err)
		return utils.SendCodeError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "code sent to " + email}}
}

func Login(db *gorm.DB, email, password string) utils.Response {
	var user User
	if err := QueryUser(db, &User{Email: email}, &user); err != nil {
		log.Print(err)
		return utils.LoginAuthError
	}
	if bcrypt.CompareHashAndPassword(user.Password, []byte(password)) != nil {
		return utils.LoginAuthError
	}
	jwt, err := utils.CreateJWT(user.ID, 24*time.Hour)
	if err != nil {
		log.Print(err)
		return utils.UserLoginError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"email": user.Email, "token": jwt}}
}

func ChangePassword(rc *redis.Client, db *gorm.DB, email, password string, code string) utils.Response {
	if !VerifyCode(rc, email, code) {
		return utils.VerificationCodeError
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return utils.ChangePasswordError
	}
	rowsAffected, err := UpdateUsers(db, &User{Email: email}, &User{Password: hashedPassword}, &[]User{})
	if err != nil {
		log.Print(err)
		return utils.ChangePasswordError
	}
	if rowsAffected == 0 {
		return utils.FindUserError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "change password success"}}
}
