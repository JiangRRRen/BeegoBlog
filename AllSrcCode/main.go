package main

import (
	"Blog/models"
	_"Blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	models.Init()
	beego.Run()
}

