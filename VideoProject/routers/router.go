package routers

import (
	"VideoProject/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//在这里面实现
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/user/login",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
	beego.Router("/", &controllers.MainController{})
}
