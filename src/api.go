package src

import (
	"fmt"
	"net/http"
	"strings"
	"traffic-forward/src/common"
	"traffic-forward/src/database"
	"traffic-forward/src/models"

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
	if m, err := models.ForwardDelete(db, uint(id)); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
	} else {
		if err := common.CloseTrans(m.Port); err != nil {
			fmt.Printf("ForwardDelete:%d,%s", m.Port, err.Error())
			ctx.JSON(http.StatusOK, models.Output{Status: 1})
			return
		}
		ctx.JSON(http.StatusOK, models.Output{Status: 1})
	}
}
func PostForwardHandle(ctx *gin.Context) {
	var body models.Forward_Req
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}

	port, err := common.ValidatePort(body.BindPort)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: fmt.Sprintf("[配置异常:%s] - %s", body.BindPort, err.Error())})
		return
	}
	body.BindPort = port
	db := database.GetDB()
	model, err := models.ForwardCreateOrUpdate(db, body)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
		return
	}

	if body.ID == 0 {
		common.RunTransferred(0, body.BindPort, body.Destination, func(use_total uint64) {
			models.ForwardUpdateUse(db, model.ID, use_total)
		})
	} else {
		if err := common.CloseTrans(model.Port); err != nil {
			ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: err.Error()})
			return
		}
		common.RunTransferred(0, body.BindPort, body.Destination, func(use_total uint64) {
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

func GetWAF(ctx *gin.Context) {
	msg, err := common.UFW_Checked()
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: msg, Data: err})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: msg})
}

func GetWAFStatus(ctx *gin.Context) {
	msg, err := common.UFW_Status()
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: msg, Data: err})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: msg})
}

func UpdateWAF(ctx *gin.Context) {
	type CMD struct {
		Cmd string `json:"cmd"`
	}
	var body CMD
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: "参数异常"})
		return
	}
	if body.Cmd == "" {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: "参数异常"})
		return
	}
	args := strings.Split(body.Cmd, " ")
	msg, err := common.UFW_Command(args)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: msg, Data: err})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: msg})
}

func DeleteWAF(ctx *gin.Context) {
	type CMD struct {
		Id string `json:"id"`
	}
	var body CMD
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: "参数异常"})
		return
	}
	if body.Id == "" {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: "参数异常"})
		return
	}
	msg, err := common.UFW_Command([]string{"-f", "delete", body.Id})
	if err != nil {
		ctx.JSON(http.StatusOK, models.Output{Status: 0, Message: msg, Data: err})
		return
	}
	ctx.JSON(http.StatusOK, models.Output{Status: 1, Data: msg})
}
