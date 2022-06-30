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

func AddUser(rc *redis.Client, db *gorm.DB, email, password, code string) utils.Response {
	if !VerifyCode(rc, email, code) {
		return utils.VerificationCodeError
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return utils.CreateUserError
	}
	rowsAffected, err := CreateUser(db, &User{Email: email, Password: hashedPassword})
	if err != nil {
		log.Print(err)
		return utils.CreateUserError
	}
	if rowsAffected == 0 {
		return utils.CreateUserError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "registration success"}}
}

func SendCode(s *utils.Sender, rc *redis.Client, email string) utils.Response {
	if err := s.SendCode(rc, email); err != nil {
		log.Print(err)
		return utils.SendCodeError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "code sent to " + email}}
}

func Login(db *gorm.DB, email, password string) utils.Response {
	var users []User
	rowsAffected, err := QueryUsers(db, &User{Email: email}, &users)
	if err != nil {
		log.Print(err)
		return utils.LoginAuthError
	}
	if rowsAffected == 0 || bcrypt.CompareHashAndPassword(users[0].Password, []byte(password)) != nil {
		return utils.LoginAuthError
	}
	jwt, err := utils.CreateJWT(users[0].Id, 24*time.Hour)
	if err != nil {
		log.Print(err)
		return utils.UserLoginError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"email": users[0].Email, "token": jwt}}
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
	rowsAffected, err := UpdateUsers(db, &User{Email: email}, &User{Password: hashedPassword})
	if err != nil {
		log.Print(err)
		return utils.ChangePasswordError
	}
	if rowsAffected == 0 {
		return utils.FindUserError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "change password success"}}
}
