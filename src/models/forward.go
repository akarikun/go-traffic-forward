package models

import (
	"time"
)

type Forward struct {
	ID          uint      `json:"id" gorm:"primarykey;autoIncrement"`
	UserID      uint      `json:"user_id"`
	Port        uint16    `json:"port" gorm:"unique"`
	Destination string    `json:"destination"`
	Ratio       float32   `json:"ratio"`
	AddDate     time.Time `json:"add_date"`
	UseTotal    uint64    `json:"use_total"`
	IsDel       int       `json:"is_del"`
}
