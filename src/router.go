package src

import (
	"TRAFforward/src/database"
	"TRAFforward/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var api_str = "/api"

func checkCookie(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == api_str+"/login.php" {
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

func RouterRegister(r *gin.Engine) {
	g := r.Group(api_str).Use(checkCookie(database.GetDB()))
	g.POST("/login.php", PostLoginHandle)
	g.GET("/forward.php", GetForwardHandle)
	g.POST("/forward.php", PostForwardHandle)
	g.POST("/forward_del.php", PostForwardDeleteHandle)
}
