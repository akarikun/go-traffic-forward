package src

import (
	"TRAFforward/src/models"
	"encoding/json"
	"time"

	"os"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Addr           string `json:"addr"`
	SqlType        int    `json:"sqlType"`
	ConnectionText string `json:"connectionText"`
}

func InitDB() (*gorm.DB, string) {
	cfg := InitConfig()

	// 初始化数据库连接
	db, err := gorm.Open(sqlite.Open(cfg.ConnectionText), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.User{},
	)

	var u []models.User
	db.Find(&u)
	if len(u) == 0 {
		u4 := uuid.New()
		user := models.User{
			UUID:          u4.String(),
			Username:      "admin",
			Nickname:      "admin",
			Password:      "123456",
			RegisterDate:  time.Now(),
			LastLoginDate: time.Now(),
			Token:         u4.String(),
			Affiliates:    "",
			Type:          1,
			IsDel:         0,
		}
		db.Create(&user)
	}

	return db, cfg.Addr
}

func InitConfig() Config {
	filename := "./config.json"

	_, err := os.Stat(filename)
	if err != nil {
		cfg := Config{
			Addr:           "127.0.0.1:8086",
			SqlType:        0, // 0:sqlite
			ConnectionText: "./data.db",
		}
		cfgData, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			panic("config.json转换异常")
		}
		os.WriteFile(filename, cfgData, 0666)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("config.json读取异常")
	}
	var cfgJson Config
	json.Unmarshal(data, &cfgJson)
	return cfgJson
}
