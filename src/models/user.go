package models

import (
	"TRAFforward/src/common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	TokenDate     time.Time `json:"token_date"`
	Affiliates    string    `json:"affiliates"`
	Type          int       `json:"type"`
	IsDel         int       `json:"is_del"`
}
type User_Resp struct {
	Username     string
	Nickname     string
	RegisterDate time.Time
	Token        string
	Affiliates   string
	Type         int
}

func UserLogin(db *gorm.DB, username, password string) User {
	var u User
	db.Where(User{Username: username, Password: password}).First(&u)
	if u.ID > 0 {
		u4 := uuid.New()
		u.Token = u4.String()
		u.TokenDate = time.Now()
		db.Save(u)
	}
	return u
}

func UserByToken(db *gorm.DB, token string) User {
	var u User
	db.Where(User{Token: token}).First(&u)
	if u.ID > 0 && u.TokenDate.AddDate(0, 0, 30).After(time.Now()) {
		return u
	}
	return User{}
}

func UserCreateOrUpdate(db *gorm.DB, req User) {
	if req.ID == 0 {
		u4 := uuid.New()
		req.UUID = u4.String()
		req.RegisterDate = time.Now()
		req.IsDel = 0
		req.Type = 1
		req.Affiliates = "1"
		req.Token = u4.String()
		db.Create(&req)
	} else {
		var u User
		db.First(&u, 1)
		if common.IsNullOrEmpty(req.Password) {
			u.Password = req.Password
		}
		db.Save(u)
	}
}
