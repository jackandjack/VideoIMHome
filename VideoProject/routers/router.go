package routers

import (
	"VideoIMHome/VideoProject/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/api/Registe", &controllers.RegisteUserController{}, "POST:Registe")
	beego.Router("/api/Login", &controllers.LoginUserController{}, "POST:Login")
}
