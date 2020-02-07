package routers

import (
	"Blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	//自动路由，比如Login()方法直接转换为admin/login
    beego.AutoRouter(&controllers.AdminController{})
}
