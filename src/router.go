package src

import (
	"net/http"
	"traffic-forward/src/database"
	"traffic-forward/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var API_BASE_URL = ""

func checkCookie(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == API_BASE_URL+"/login.do" {
			ctx.Next()
			return
		}
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(http.StatusOK, models.Output{Status: -1, Message: err.Error()})
			ctx.Abort()
			return
		}
		u := models.UserByToken(db, token)
		if u.ID == 0 {
			ctx.JSON(http.StatusOK, models.Output{Status: -1, Message: "用户信息异常"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func RouterRegister(r *gin.Engine, cfg database.Config) {
	API_BASE_URL = cfg.APIBaseUrl
	g := r.Group(cfg.APIBaseUrl).Use(checkCookie(database.GetDB()))
	g.POST("/login.do", PostLoginHandle)
	g.GET("/forward.do", GetForwardHandle)
	g.POST("/forward.do", PostForwardHandle)
	g.POST("/forward_del.do", PostForwardDeleteHandle)
	g.GET("/waf.do", GetWAF)
	g.GET("/waf_status.do", GetWAFStatus)
	g.POST("/waf_update.do", UpdateWAF)
	g.POST("waf_delete.do", DeleteWAF)
}
