package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func Init(){
	//从配置文件中获取连接参数
	host:=beego.AppConfig.String("mysqlHost")
	port:=beego.AppConfig.String("mysqlPort")
	user:=beego.AppConfig.String("mysqlUser")
	password:=beego.AppConfig.String("mysqlPassword")
	dbName:=beego.AppConfig.String("mysqlDb")

	//Data Source Name=DSN
	connectInfo:=[]string{user,":",password,"@tcp(",host,":",port,")/", dbName, "?charset=utf8"}
	DSN:=strings.Join(connectInfo,"")
	orm.RegisterDataBase("default","mysql",DSN)
	orm.RegisterModel(new(User),new(Config),new(Category),new(Post))

}

func GetTableName(str string) string{
	return beego.AppConfig.String("mysqlPrefix")+str
}
