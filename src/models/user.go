package models

import (
	"TRAFforward/src/common"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UUID          string    `json:"uuid" gorm:"unique;not null"`
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

type User_Query struct {
	*Query
}

func newAff() string {
	return common.Random(12)
}

func UserCreateAdmin(db *gorm.DB) {
	var m []User
	db.Find(&m)
	if len(m) == 0 {
		user := User{
			UUID:          common.UUID(),
			Username:      "admin",
			Nickname:      "admin",
			Password:      common.MD5("123456"),
			RegisterDate:  time.Now(),
			LastLoginDate: time.Now(),
			Token:         common.UUID(),
			Affiliates:    newAff(),
			Type:          1,
			IsDel:         0,
		}
		db.Create(&user)
	}
}

func UserLogin(db *gorm.DB, username, password string) User {
	var u User
	db.Where(User{Username: username, Password: common.MD5(password)}).First(&u)
	if u.ID > 0 {
		token := common.UUID()
		tokenDate := time.Now().AddDate(0, 0, 30) //有效期30天
		if u.TokenDate.Before(time.Now()) {       //过期
			u.Token = token
			u.TokenDate = tokenDate
			db.Save(u)
		}
	}
	return u
}

func UserByToken(db *gorm.DB, token string) User {
	var u User
	db.Where(User{Token: token}).First(&u)
	if u.ID > 0 && u.TokenDate.After(time.Now()) {
		return u
	}
	return User{}
}

func UserCreateOrUpdate(db *gorm.DB, req User) {
	if req.ID == 0 {
		req.UUID = common.UUID()
		req.RegisterDate = time.Now()
		req.IsDel = 0
		req.Type = 1
		req.Affiliates = newAff()
		req.Token = common.UUID()
		db.Create(&req)
	} else {
		var m User
		db.First(&m, 1)
		if common.IsNullOrEmpty(req.Password) {
			m.Password = common.MD5(req.Password)
		}
		db.Save(m)
	}
}
