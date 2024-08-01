package models

import (
	"time"

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

type Forward_Req struct {
	Req
}

func portChecked(db *gorm.DB, port uint16, id uint) bool {
	var count int64
	db.Where(Forward{Port: port}).Not(Forward{ID: id}).Count(&count)
	return count > 0
}
func ForwardGetList(db *gorm.DB, req Forward_Req) []Forward {
	var m []Forward
	db.Limit(req.PageSize).Offset((req.PageIndex - 1) * req.PageSize).Find(&m)
	return m
}
func ForwardCreateOrUpdate(db *gorm.DB, req Forward) {
	if portChecked(db, req.Port, req.ID) {
		if req.ID == 0 {
			req.Ratio = 1
			req.AddDate = time.Now()
			req.UseTotal = 0
			req.IsDel = 0
			db.Create(&req)
		} else {
			var m Forward
			db.First(&m, 1)
			m.BindPort = req.BindPort
			m.Destination = req.Destination
			m.Port = req.Port
			db.Save(m)
		}
	}
}
