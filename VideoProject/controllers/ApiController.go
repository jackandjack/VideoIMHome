package controllers

import (
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Login() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	username := this.GetString("username")
	password := this.GetString("password")
	if len(username) == 0 {
		//提示账号错误

	}
	if len(password) == 0 {
		//提示密码错误
	}
	this.RetData(resp)
}

func (this *UserController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
