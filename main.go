package main

import (
	"TRAFforward/src"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	go src.Transferred(":8085", "127.0.0.1:57890", func(use uint64) {
		log.Printf("8085: %d", use)
	})
	go src.Transferred(":8086", "127.0.0.1:57890", func(use uint64) {
		log.Printf("8086: %d", use)
	})
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
