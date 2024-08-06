package src

import (
	"TRAFforward/src/common"
	"TRAFforward/src/database"
	"TRAFforward/src/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostForwardDeleteHandle(ctx *gin.Context) {
	var body map[string]int
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Message: err.Error()})
		return
	}
	id := body["id"]
	db := database.GetDB()
	m := models.ForwardDelete(db, uint(id))
	if err := common.CloseTrans(m.Port); err != nil {
		// ctx.JSON(http.StatusOK, models.Output{Message: err.Error()})
		fmt.Printf("ForwardDelete:%d,%s", m.Port, err.Error())
		ctx.JSON(http.StatusOK, models.Output{Status: 1})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1})
}
func PostForwardHandle(ctx *gin.Context) {
	var body models.Forward_Req
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}

	// var bool,
	ok, err := common.ValidatePort(body.BindPort)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: fmt.Sprintf("端口[%s]已被占用", body.BindPort)})
		return
	}
	db := database.GetDB()
	model, err := models.ForwardCreateOrUpdate(db, body)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}
	if body.ID == 0 {
		go common.RunTransferred(0, 10, body.BindPort, body.Destination, func(use_total uint64) {
			models.ForwardUpdateUse(db, model.ID, use_total)
		})
	} else {
		if err := common.CloseTrans(model.Port); err != nil {
			ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
			return
		}
		go common.RunTransferred(model.UseTotal, 10, body.BindPort, body.Destination, func(use_total uint64) {
			models.ForwardUpdateUse(db, model.ID, use_total)
		})
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
		ctx.JSON(http.StatusOK, models.Output{Message: err.Error()})
		return
	}
	username, ok1 := body["username"]
	password, ok2 := body["password"]
	if !(ok1 && ok2) {
		ctx.JSON(http.StatusOK, models.Output{Message: "缺少参数"})
		return
	}
	db := database.GetDB()
	u := models.UserLogin(db, username, password)
	if u.ID == 0 {
		ctx.JSON(http.StatusOK, models.Output{Message: "用户名或密码不正确"})
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
