package main

import (
	"TRAFforward/src"
	"TRAFforward/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	go src.Transferred(0, "127.0.0.1:30003", "127.0.0.1:59992", func(use uint64) {
		//log.Printf("59992: %d", use)
	})
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})

	db, addr := src.InitDB()

	// 创建用户
	r.POST("/users", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&user)
		c.JSON(200, user)
	})

	// 获取所有用户
	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		db.Find(&users)
		c.JSON(200, users)
	})
	r.Run(addr)
}
