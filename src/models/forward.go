package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Forward struct {
	ID          uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UserID      uint      `json:"user_id"`
	Port        uint16    `json:"port" gorm:"unique"`
	BindPort    string    `json:"bind_port"`
	Destination string    `json:"destination"`
	Ratio       float32   `json:"ratio"`
	AddDate     time.Time `json:"add_date"`
	UseTotal    uint64    `json:"use_total"`
	IsDel       int       `json:"is_del"`
}

type Forward_Query struct {
	Query
}
type Forward_Req struct {
	ID          uint   `json:"id" gorm:"primarykey;autoIncrement"`
	UserID      uint   `json:"user_id"`
	BindPort    string `json:"bind_port"`
	Destination string `json:"destination"`
}

func getPort(db *gorm.DB, bind_port string) (uint16, error) {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(bind_port)
	_port, err := strconv.Atoi(match)
	if err != nil {
		return 0, errors.New("端口异常")
	}
	port := uint16(_port)
	var count int64
	db.Where(&Forward{Port: port, IsDel: 0}).Count(&count)
	if count > 0 {
		return 0, errors.New("端口已占用")
	}
	return port, nil
}
func ForwardGetList(db *gorm.DB, query Forward_Query) []Forward {
	var m []Forward
	db.Where("is_del", 0).Limit(query.PageSize).Offset((query.PageIndex - 1) * query.PageSize).Order("id desc").Find(&m)
	return m
}
func ForwardCreateOrUpdate(db *gorm.DB, req Forward_Req) error {
	port, err := getPort(db, req.BindPort)
	if err != nil {
		return err
	}
	if req.ID == 0 {
		m := Forward{
			Port:     port,
			Ratio:    1,
			AddDate:  time.Now(),
			UseTotal: 0,
			IsDel:    0,
		}
		copier.Copy(&m, &req)
		result := db.Create(&m)
		fmt.Println(result)
	} else {
		var m Forward
		copier.Copy(&m, &req)
		db.Updates(m)
	}
	return nil
}
func ForwardDelete(db *gorm.DB, id uint) {
	db.Where(Forward{ID: id}).Updates(Forward{IsDel: 1})
}
