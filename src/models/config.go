package models

type Config struct {
	ID     uint   `json:"id" gorm:"primarykey"`
	Listen string `json:"listen"`
	Port   int    `json:"port"`
}
