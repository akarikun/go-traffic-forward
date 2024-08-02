package main

import (
	"TRAFforward/src"
	"TRAFforward/src/database"
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
	r := gin.Default()
	templ := template.Must(template.New("").ParseFS(templatesEmbed, "www/*.html"))
	r.SetHTMLTemplate(templ)
	jsFS, _ := fs.Sub(JSEmbed, "www/js")
	r.StaticFS("/js", http.FS(jsFS))
	cssFS, _ := fs.Sub(CSSEmbed, "www/css")
	r.StaticFS("/css", http.FS(cssFS))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	cfg := database.InitConfig()
	db := database.InitDB(cfg)
	db.AutoMigrate(
		&models.User{},
		&models.Forward{},
	)
	models.UserCreateAdmin(db)
	src.RouterRegister(r)
	r.Run(cfg.Addr)
}
