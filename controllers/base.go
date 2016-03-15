package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"

	"github.com/shwinpiocess/cc/models"
	"github.com/shwinpiocess/cc/utils"
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *models.User
	userId         int
	userName       string
	pageSize       int
}

func (this *BaseController) Prepare() {
	this.pageSize = 20
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)
	this.auth()

	this.Data["curRoute"] = this.controllerName + "." + this.actionName
	this.Data["curController"] = this.controllerName
	this.Data["curAction"] = this.actionName
	this.Data["loginUserId"] = this.userId
	this.Data["loginUserName"] = this.userName
}

func (this *BaseController) getClientIP() string {
	return strings.Split(this.Ctx.Request.RemoteAddr, ":")[0]
}

func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
	this.StopRun()
}

func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}

func (this *BaseController) jsonResult(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) auth() {
	fmt.Println("auth*****************************************")
	arrs := strings.Split(this.Ctx.GetCookie("auth"), "|")
	fmt.Println("arrs=", arrs)
	if len(arrs) == 2 {
		idstr, password := arrs[0], arrs[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := models.GetUserById(userId)
			fmt.Println("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy", user, err)
			if err == nil && password == utils.Md5([]byte(this.getClientIP()+"|"+user.Password+user.Salt)) {
				this.userId = user.Id
				this.userName = user.UserName
				this.user = user
			}
		}
	}

	if this.userId == 0 && (this.controllerName != "main" || (this.controllerName == "main" && this.actionName != "logout" && this.actionName != "login")) {
		this.redirect(beego.URLFor("MainController.Login"))
	}
}
