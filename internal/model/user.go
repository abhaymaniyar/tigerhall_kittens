package model

type User struct {
	Username string `gorm:"unique"`
	Password string
	Email    string `gorm:"unique"`
}
