package database

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	Addr           string `json:"addr"`
	SqlType        int    `json:"sqlType"`
	ConnectionText string `json:"connectionText"`
}

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB(cfg Config) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(cfg.ConnectionText), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,  // 慢 SQL 阈值
				LogLevel:                  logger.Error, // 日志级别
				IgnoreRecordNotFoundError: true,         // 忽略 record not found 错误
				Colorful:                  false,        // 禁用彩色打印
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
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
