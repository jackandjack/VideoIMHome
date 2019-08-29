package routers

import (
	"VideoIMHome/VideoProject/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/api", &controllers.UserController{}, "*:Login")
}
