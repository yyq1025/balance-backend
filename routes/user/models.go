package user

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
)

type User struct {
	ID       int
	Email    string
	Password []byte
}

type Sender struct {
	email, password, smtpHost, tlsPort string
	auth                               smtp.Auth
}

func newSender() *Sender {
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	tlsPort := os.Getenv("MAIL_PORT")
	auth := smtp.PlainAuth("", email, password, smtpHost)
	return &Sender{email, password, smtpHost, tlsPort, auth}
}

func (s *Sender) sendCode(rc_cache *cache.Cache, to string) error {
	subject := "Subject: Your OTP Code.\n\n"
	code := fmt.Sprintf("%06d", rand.Intn(1e6))
	ending := "\n\nCode expires in 30 min."
	if err := smtp.SendMail(s.smtpHost+":"+s.tlsPort, s.auth, s.email, []string{to}, []byte(subject+code+ending)); err != nil {
		return err
	}
	return rc_cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   fmt.Sprintf("code:%s", to),
		Value: code,
		TTL:   time.Duration(30 * time.Minute),
	})
}
