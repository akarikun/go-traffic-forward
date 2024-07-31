package src

import (
	"TRAFforward/src/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/robfig/cron"
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
		// Logger: logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		// 	logger.Config{
		// 		SlowThreshold:             time.Second, // 慢 SQL 阈值
		// 		LogLevel:                  logger.Warn, // 日志级别
		// 		IgnoreRecordNotFoundError: true,        // 忽略 record not found 错误
		// 		Colorful:                  false,       // 禁用彩色打印
		// 	},
		// ),
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

func FormatUse(use uint64) string {
	const (
		KB = 1 << (10 * (iota + 1))
		MB
		GB
		TB
	)

	switch {
	case use >= TB:
		return fmt.Sprintf("%.2f TB", float64(use)/TB)
	case use >= GB:
		return fmt.Sprintf("%.2f GB", float64(use)/GB)
	case use >= MB:
		return fmt.Sprintf("%.2f MB", float64(use)/MB)
	case use >= KB:
		return fmt.Sprintf("%.2f KB", float64(use)/KB)
	default:
		return fmt.Sprintf("%d B", use)
	}
}

func RunTransferred(value uint64, sourcePort string, destinationAddress string) {
	var m sync.Mutex
	var use uint64 = 1 //值为0时会重置使用量
	go Transferred(value, sourcePort, destinationAddress, func(_use uint64, _cur int) uint {
		//log.Printf("30003: %d, %s", _use, FormatUse(_use))
		if use == 0 {
			// log.Printf("reset use %d,%s", _use, FormatUse(_use))
			use = 1
			return 1
		} else {
			use = _use
		}
		return 100 //0停止 1重新统计
	})
	log.Printf("任务执行时间：%s", time.Now())
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		m.Lock()
		defer m.Unlock()
		if use > 0 {
			//log.Printf("任务执行时间：%s,%d,%s", time.Now(), use, FormatUse(use))
			//需要统计流量
		}
		use = 0 //重置状态
	})
	c.Start()
}
