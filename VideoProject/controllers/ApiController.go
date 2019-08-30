package controllers

import (
	"VideoIMHome/VideoProject/models"
	_ "VideoIMHome/VideoProject/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	SecretKey = "VideoProgress_2016"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Regist() {

	username := this.GetString("username")

	password := this.GetString("password")

}

func GetToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(SecretKey))
	fmt.Println("Token=", tokenString)
	return tokenString
}

func (this *UserController) Login() {

	user := models.User{}
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	username := this.GetString("username")
	password := this.GetString("password")
	if len(username) == 0 {
		//提示账号错误
		resp["errno"] = 10001
		resp["errmsg"] = "请输入用户名"
		this.RetData(resp)
		return
	}
	if len(password) == 0 {
		//提示密码错误
		resp["errno"] = 10002
		resp["errmsg"] = "请输入用户密码"
		this.RetData(resp)
		return
	}
	user.Username = username
	user.Password = password
	newuser, suess := models.Queryall(&user)
	if suess {
		resp["errno"] = 10000
		str, _ := json.Marshal(newuser)
		fmt.Printf("%s\n", str)
		resp["data"] = fmt.Sprintf("%s", str)
		resp["errmsg"] = "登录成功"
	} else {
		issuess := models.AddUser(user)
		if issuess {
			resp["errno"] = 10000
			dataJ, _ := json.Marshal(newuser)
			resp["data"] = dataJ
			resp["errmsg"] = "登录成功"
		} else {
			resp["errno"] = 10004
			resp["errmsg"] = "登录失败"
		}
	}
	this.RetData(resp)
}

func (this *UserController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
