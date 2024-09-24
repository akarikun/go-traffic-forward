package models

import (
	"TRAFforward/src/common"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Forward struct {
	ID          uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UserID      uint      `json:"user_id"`
	Port        uint16    `json:"port"`
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
	port, _, err := common.GetPort(bind_port)
	if err != nil {
		return 0, err
	}
	var count int64
	db.Model(&Forward{}).Where("port=? and is_del=0", port).Count(&count)
	if count > 0 {
		return 0, errors.New("端口已占用")
	}
	return port, nil
}
func ForwardGetList(db *gorm.DB, query Forward_Query) []Forward {
	var m []Forward
	db.Model(Forward{}).Where("is_del", 0).Limit(query.PageSize).Offset((query.PageIndex - 1) * query.PageSize).Order("id desc").Find(&m)
	return m
}

func ForwardGetPortList(db *gorm.DB) []Forward {
	var m []Forward
	db.Model(Forward{}).Where("is_del", 0).Find(&m)
	return m
}

func ForwardCreateOrUpdate(db *gorm.DB, req Forward_Req) (Forward, error) {
	port, err := getPort(db, req.BindPort)
	if err != nil {
		return Forward{}, err
	}
	if req.ID == 0 {
		m := Forward{
			BindPort: req.BindPort,
			Port:     port,
			Ratio:    1,
			AddDate:  time.Now(),
			UseTotal: 0,
			IsDel:    0,
		}
		copier.Copy(&m, &req)
		if fmt.Sprintf("%d", port) == m.BindPort {
			m.BindPort = fmt.Sprintf("127.0.0.1:%d", port)
		}
		result := db.Create(&m)
		return m, result.Error
	} else {
		var m Forward
		copier.Copy(&m, &req)
		if fmt.Sprintf("%d", port) == m.BindPort {
			m.BindPort = fmt.Sprintf("127.0.0.1:%d", port)
		}
		result := db.Updates(m)
		return m, result.Error
	}
}
func ForwardDelete(db *gorm.DB, id uint) (Forward, error) {
	var m Forward
	if err := db.First(&m, id).Error; err != nil {
		return m, err // 返回错误
	}
	// m.IsDel = 1
	// if err := db.Save(&m).Error; err != nil {
	// 	return m, err // 返回错误
	// }
	// return m, nil
	db.Delete(&Forward{}, id)
	return m, nil
}
func ForwardUpdateUse(db *gorm.DB, id uint, use_total uint64) {
	var m Forward
	db.Where(Forward{ID: id}).Find(&m)
	m.UseTotal += use_total
	db.Save(m)
	// db.Where(Forward{ID: id}).Updates(Forward{UseTotal: use_total})
}
