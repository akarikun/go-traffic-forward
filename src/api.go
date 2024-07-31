package src

import (
	"TRAFforward/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Api(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	g := r.Group("/api")
	g.POST("/login.php", func(ctx *gin.Context) {
		var jsonData map[string]string
		if err := ctx.BindJSON(&jsonData); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: err.Error()})
			return
		}
		username, ok1 := jsonData["username"]
		password, ok2 := jsonData["password"]
		if !(ok1 && ok2) {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: "缺少参数"})
			return
		}
		u := models.UserLogin(db, username, password)
		if u.ID == 0 {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: "用户名或密码不正确"})
			return
		}
		ctx.SetCookie("token", u.Token, 60*60*24*30, "/", ctx.Request.Host, false, true)
		resp := models.User_Resp{
			Username:     u.Username,
			Nickname:     u.Nickname,
			RegisterDate: u.RegisterDate,
			Token:        u.Token,
			Affiliates:   u.Affiliates,
			Type:         u.Type,
		}
		ctx.JSON(http.StatusOK, models.Resp{Status: 1, Data: resp})
	})
	g.POST("/post_forward.php", func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(http.StatusNotFound, models.Resp{})
			return
		}
		u := models.UserByToken(db, token)
		if u.ID == 0 {
			ctx.JSON(http.StatusNotFound, models.Resp{})
			return
		}
		ctx.JSON(http.StatusOK, models.Resp{Message: "缺少参数"})
	})
}
