package main

import (
	_ "VideoIMHome/VideoProject/controllers"
	_ "VideoIMHome/VideoProject/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()

}
