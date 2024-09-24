package models

import "time"

type WAF struct {
	ID              uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UserID          uint      `json:"user_id"`
	IP              string    `json:"ip"`   //对方IP信息
	Port            uint      `json:"port"` //服务端端口
	RequestDateTime time.Time `json:"request_date_time"`
}
