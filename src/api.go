package src

import (
	"TRAFforward/src/database"
	"TRAFforward/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostForwardDeleteHandle(ctx *gin.Context) {
	var body map[string]int
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Output{Message: err.Error()})
		return
	}
	id := body["id"]
	db := database.GetDB()
	models.ForwardDelete(db, uint(id))
}
func PostForwardHandle(ctx *gin.Context) {
	var body models.Forward_Req
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}
	db := database.GetDB()
	if err := models.ForwardCreateOrUpdate(db, body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1})
}

func GetForwardHandle(ctx *gin.Context) {
	var body models.Forward_Query
	if err := ctx.BindQuery(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: "参数异常"})
		return
	}
	db := database.GetDB()
	list := models.ForwardGetList(db, body)
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: list})
}

func PostLoginHandle(ctx *gin.Context) {
	var body map[string]string
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Output{Message: err.Error()})
		return
	}
	username, ok1 := body["username"]
	password, ok2 := body["password"]
	if !(ok1 && ok2) {
		ctx.JSON(http.StatusBadRequest, models.Output{Message: "缺少参数"})
		return
	}
	db := database.GetDB()
	u := models.UserLogin(db, username, password)
	if u.ID == 0 {
		ctx.JSON(http.StatusBadRequest, models.Output{Message: "用户名或密码不正确"})
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
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: resp})
}
