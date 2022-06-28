package user

type User struct {
	Id       int `gorm:"autoIncrement"`
	Email    string
	Password []byte
}
