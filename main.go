package main

import (
	"TRAFforward/src"
	"TRAFforward/src/models"
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed www/js
	JSEmbed embed.FS

	//go:embed www/css
	CSSEmbed embed.FS

	//go:embed www/*
	templatesEmbed embed.FS
)

func main() {
	go src.Transferred(0, "127.0.0.1:30003", "127.0.0.1:57890", func(use uint64) {
		//log.Printf("59992: %d", use)
	})
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	templ := template.Must(template.New("").ParseFS(templatesEmbed, "www/*.html"))
	r.SetHTMLTemplate(templ)
	jsFS, _ := fs.Sub(JSEmbed, "www/js")
	r.StaticFS("/js", http.FS(jsFS))
	cssFS, _ := fs.Sub(CSSEmbed, "www/css")
	r.StaticFS("/css", http.FS(cssFS))

	r.GET("/", func(c *gin.Context) {
		// c.String(http.StatusOK, "hello World!")
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
	})
	db, addr := src.InitDB()
	r.POST("/users", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&user)
		c.JSON(200, user)
	})
	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		db.Find(&users)
		c.JSON(200, users)
	})
	r.Run(addr)
}
