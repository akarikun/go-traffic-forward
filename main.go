package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"traffic-forward/src"
	"traffic-forward/src/common"
	"traffic-forward/src/database"
	"traffic-forward/src/models"

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

func InitForward() {
	db := database.GetDB()
	list := models.ForwardGetPortList(db)
	for _, v := range list {
		port, err := common.ValidatePort(v.BindPort)
		if err != nil {
			fmt.Printf("[配置异常:%s] - %s", port, err.Error())
			continue
			// return
		}
		common.RunTransferred(0, port, v.Destination, func(use_total uint64) {
			models.ForwardUpdateUse(db, v.ID, use_total)
		})
	}
}

func main() {
	cfg := database.InitConfig()
	if !cfg.Debugger {
		gin.SetMode(gin.ReleaseMode)
	}
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
	db := database.InitDB(cfg)
	db.AutoMigrate(
		&models.User{},
		&models.Forward{},
	)
	models.UserCreateAdmin(db)
	src.RouterRegister(r)
	InitForward()
	r.Run(cfg.Addr)
}
