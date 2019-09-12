package controllers

import (
	"VideoIMHome/VideoProject/models"
	_ "VideoIMHome/VideoProject/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/dgrijalva/jwt-go/request"
	"strconv"
	"time"
	_ "time"
)

type OprationType int32

const (
	OprationTypeIcon         OprationType = 0
	OprationTypeNickName     OprationType = 1
	OprationTypePwd          OprationType = 2
	OprationTypeSex          OprationType = 3
	OprationTypeIntroduction OprationType = 4
	OprationTypeBirthday     OprationType = 5
)

const (
	SecretKey = "VideoProgress_2016"
)

func Validate(isLogin bool, this *RegisteUserController, c *LoginUserController) bool {
	if !isLogin {
		phone := this.GetString("phone")
		sendType := this.GetString("sendType")
		resp := make(map[string]interface{})
		if len(phone) == 0 {
			resp["code"] = 1002
			resp["data"] = ""
			resp["msg"] = "请输入电话号码"
			resp["info_code"] = 1002
			this.RetData(resp)
			return false
		}
		if len(sendType) == 0 {
			resp["code"] = 1003
			resp["data"] = ""
			resp["msg"] = "没有的类型"
			resp["info_code"] = 1003
			this.RetData(resp)
			return false
		}
		return true
	} else {
		phone := c.GetString("phone")
		resp := make(map[string]interface{})
		if len(phone) == 0 {
			resp["code"] = 1002
			resp["data"] = ""
			resp["msg"] = "请输入电话号码"
			resp["info_code"] = 1002
			c.RetData(resp)
			return false
		}
	}
	return true
}

type RegisteUserController struct {
	beego.Controller
}

type LoginUserController struct {
	beego.Controller
}

type UpdateUserController struct {
	beego.Controller
}

func check(optype OprationType, token string, text string, this *UpdateUserController) bool {

	resp := make(map[string]interface{})
	//判断操作是否存在
	if optype != OprationTypeNickName && optype != OprationTypeIcon && optype != OprationTypeBirthday && optype != OprationTypeIntroduction && optype != OprationTypeSex && optype != OprationTypePwd {
		resp["code"] = 1101
		resp["data"] = ""
		resp["msg"] = "该操作不存在"
		resp["info_code"] = 1101
		this.RetData(resp)
		return false
	}
	//判断Token 是否存在
	if models.CheckToken(token) == false {
		resp["code"] = 1111
		resp["data"] = ""
		resp["msg"] = "Token失效"
		resp["info_code"] = 1111
		this.RetData(resp)
		return false
	}
	//判断文本是否未空
	if len(text) == 0 {
		resp["code"] = 1102
		resp["data"] = ""
		resp["msg"] = "文本不能为空"
		resp["info_code"] = 1102
		this.RetData(resp)
		return false
	}
	return true
}

func (this *UpdateUserController) Update() {
	resp := make(map[string]interface{})
	optype, _ := strconv.Atoi(this.GetString("type"))
	token := this.GetString("token")
	text := this.GetString("text")
	if check(OprationType(optype), token, text, this) == false {
		return
	}
	if changeUserDate(OprationType(optype), token, text) {
		resp["code"] = 1000
		resp["data"] = ""
		resp["msg"] = "更新成功"
		resp["info_code"] = 1000
	} else {
		resp["code"] = 1100
		resp["data"] = ""
		resp["msg"] = "更新失败"
		resp["info_code"] = 1100
	}
	this.RetData(resp)
}

func changeUserDate(op OprationType, token string, text string) bool {
	user := models.User{}
	switch op {
	case OprationTypeIcon:
		return models.UpdateUser(user, "Icon", text, token, "token")
	case OprationTypeNickName:
		return models.UpdateUser(user, "nickname", text, token, "token")
	case OprationTypeSex:
		return models.UpdateUser(user, "sex", text, token, "token")
	case OprationTypeIntroduction:
		return models.UpdateUser(user, "introduction", text, token, "token")
	case OprationTypeBirthday:
		return models.UpdateUser(user, "birthday", text, token, "token")
	}
	return false
}

func (this *RegisteUserController) Registe() {
	resp := make(map[string]interface{})
	user := models.User{}
	user.Phone = this.GetString("phone")
	user.Type = this.GetString("sendType")
	if Validate(false, this, nil) {
		_, isSuces := models.Regist(user)
		if isSuces {
			resp["code"] = 1000
			//str, _ := json.Marshal(usernew)
			resp["data"] = ""
			resp["msg"] = "注册成功"
			resp["info_code"] = 1000
		} else {
			resp["code"] = 1001
			//str, _ := json.Marshal(usernew)
			resp["data"] = ""
			resp["msg"] = "注册失败"
			resp["info_code"] = 1001
		}
	}
	this.RetData(resp)
}
func (this *LoginUserController) Login() {
	user := models.User{}
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	user.Phone = this.GetString("phone")
	codeStr := this.GetString("code")
	if len(codeStr) == 0 {
		resp["code"] = 1005
		resp["data"] = ""
		resp["msg"] = "请输入手机验证码"
		resp["info_code"] = 1005
		this.RetData(resp)
		return
	}
	if Validate(true, nil, this) {
		newuser, sucess := models.Findphone(user.Phone)
		if sucess {
			resp["code"] = 1000
			str, _ := json.Marshal(newuser)
			resp["data"] = fmt.Sprintf("%s", str)
			resp["msg"] = "登录成功"
			resp["info_code"] = 1000
			//更新时间戳
			t := time.Now()
			fmt.Println(t)
			t1 := t.Unix()
			times := fmt.Sprintf("%d", t1)
			isSucess := models.UpdateUser(newuser, "Token", times, newuser.Phone, "phone")
			models.UpdateUser(newuser, "Token", models.GetToken(user.Phone), newuser.Phone, "phone")
			if isSucess {
				fmt.Println("更新成功")
			} else {
				fmt.Println("更新失败")
			}
		} else {
			resp["code"] = 1006
			resp["msg"] = "登录失败"
			resp["data"] = ""
			resp["info_code"] = 1006
		}
	}
	this.RetData(resp)
}
func (this *UpdateUserController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *RegisteUserController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *LoginUserController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
