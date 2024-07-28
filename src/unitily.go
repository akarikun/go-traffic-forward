package src

import (
	"TRAFforward/src/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, string) {
	// 初始化数据库连接
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Config{},
	)
	var cfg []models.Config
	db.Find(&cfg)
	var addr string
	if len(cfg) == 0 {
		cfgVal := models.Config{
			ID:     1,
			Listen: "127.0.0.1",
			Port:   8086,
		}
		db.Create(&cfgVal)
		addr = fmt.Sprintf("%s:%d", cfgVal.Listen, cfgVal.Port)
	} else {
		addr = fmt.Sprintf("%s:%d", cfg[0].Listen, cfg[0].Port)
	}
	return db, addr
}
