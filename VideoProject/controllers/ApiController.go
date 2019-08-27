package controllers

import (
	"VideoProject/models"
	"github.com/astaxie/beego"
)

type UserController struct {

	beego.Controller

}


func (u *UserController) Login(){
	var user models.User
	username := u.GetString("username")
	password := u.GetString("password")
    suess :=models.Login(user,username,password)
	if suess {
		u.Data["json"] = "login success"
	}else{
		u.Data["json"] = "user not exist"
	}
	u.ServeJSON()

}




