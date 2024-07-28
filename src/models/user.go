package models

import (
	"time"
)

type User struct {
	ID            uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UUID          string    `json:"uuid" gorm:"unique"`
	Username      string    `json:"username" gorm:"unique"`
	Password      string    `json:"password"`
	Nickname      string    `json:"nickname"`
	RegisterDate  time.Time `json:"register_date"`
	LastLoginDate time.Time `json:"last_login_date"`
	Token         string    `json:"token" gorm:"unique"`
	Affiliates    string    `json:"affiliates"`
	Type          int       `json:"type"`
	IsDel         int       `json:"is_del"`
}
