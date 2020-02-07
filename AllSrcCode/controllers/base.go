package controllers

import (
	_"Blog/models"
	"beego-blog-blog-dev2/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

type baseController struct{
	beego.Controller
	o orm.Ormer
	controllerName string
	actionName string
}
//后台登录
//不允许在未登录的情况下访问我们的controller，所以需要判断是否已经登录
func (p *baseController) Prepare(){
	controllerName, actionName:=p.GetControllerAndAction()
	// 用strings.ToLower  转换小写提取0-5 之前的字母admin
	p.controllerName=strings.ToLower(controllerName[0:len(controllerName)-10])
	//Login 转换为小写 login
	p.actionName=strings.ToLower(actionName)

	//新建一个ORM，之前已经在initMode中注册了mysql
	p.o=orm.NewOrm()
	//跳转条件：是admin，当前动作不是login，且用户未登录（在session中没有记录）
	if strings.ToLower(p.controllerName) == "admin" && strings.ToLower(p.actionName) != "login" {
		if p.GetSession("user") == nil {
			p.History("您还未登录，请登录！", "/admin/login")
			//p.Ctx.WriteString(p.controllerName +"==="+ p.actionName)
		}
	}

	//初始化前台页面相关元素
	if strings.ToLower(p.controllerName)=="blog"{
		p.Data["actionName"]=strings.ToLower(actionName)
		var res []*models.Config
		p.o.QueryTable(new(models.Config).TableName()).All(&res)

		configs:=make(map[string]string)
		for _,v:=range res{
			configs[v.Name]=v.Value
		}
		p.Data["config"]=configs
	}

}

func (p *baseController) History(msg string, url string){
	if url==""{
		p.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		p.StopRun()
	}else{
		p.Redirect(url,302)
	}
}

//获取用户IP地址
func (p *baseController) getClientIp() string{
	s:=strings.Split(p.Ctx.Request.RemoteAddr, ":")
	return s[0]
}