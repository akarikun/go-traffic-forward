package main

import (
	"TRAFforward/src"
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
	// go src.RunTransferred()

	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	templ := template.Must(template.New("").ParseFS(templatesEmbed, "www/*.html"))
	r.SetHTMLTemplate(templ)
	jsFS, _ := fs.Sub(JSEmbed, "www/js")
	r.StaticFS("/js", http.FS(jsFS))
	cssFS, _ := fs.Sub(CSSEmbed, "www/css")
	r.StaticFS("/css", http.FS(cssFS))

	db, addr := src.InitDB()
	src.Api(r, db)
	r.Run(addr)
}
