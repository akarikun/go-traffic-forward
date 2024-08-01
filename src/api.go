package src

import (
	"TRAFforward/src/models"
	"encoding/json"
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
			ctx.JSON(http.StatusNotFound, models.Resp{})
			ctx.Abort()
			return
		}
		u := models.UserByToken(db, token)
		if u.ID == 0 {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: "接口异常"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func Api(r *gin.Engine, db *gorm.DB) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	g := r.Group(api_str).Use(checkCookie(db))
	g.POST("/login.php", func(ctx *gin.Context) {
		var body map[string]string
		if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: err.Error()})
			return
		}
		username, ok1 := body["username"]
		password, ok2 := body["password"]
		if !(ok1 && ok2) {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: "缺少参数"})
			return
		}
		u := models.UserLogin(db, username, password)
		if u.ID == 0 {
			ctx.JSON(http.StatusBadRequest, models.Resp{Message: "用户名或密码不正确"})
			return
		}
		//ctx.SetCookie("token", u.Token, 60*60*24*30, "/", ctx.Request.Host, false, true)
		ctx.SetCookie("token", u.Token, 60*60*24*30, "/", "localhost:8080", false, false)
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

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
	g.GET("/forward.php", func(ctx *gin.Context) {
		rawData, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusOK, models.Resp{Status: 0, Message: "参数异常"})
			return
		}
		var body models.Forward_Req
		if err := json.Unmarshal(rawData, &body); err != nil {
			ctx.JSON(http.StatusOK, models.Resp{Status: 0, Message: "参数异常"})
			return
		}
		list := models.ForwardGetList(db, body)
		ctx.JSON(http.StatusOK, models.Resp{Status: 1, Data: list})
	})
	g.POST("/forward.php", func(ctx *gin.Context) {
		var body models.Forward
		ctx.ShouldBindBodyWithJSON(&body)
		models.ForwardCreateOrUpdate(db, body)
		ctx.JSON(http.StatusOK, models.Resp{Status: 1})
	})
}
