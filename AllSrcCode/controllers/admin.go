package controllers

import (
	"Blog/models"
	"Blog/util"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AdminController struct{
	baseController
}


//后台用户登录
func (c *AdminController) Login(){

	if c.Ctx.Request.Method=="POST"{

		//从前台的请求中获取账号密码信息
		username:= c.GetString("username")
		password:= c.GetString("password")

		// 根据前面的信息从数据库中查出 username 的数据
		user := models.User{Username: username}
		c.o.Read(&user, "username")

		if user.Password == "" {
			c.History("账号不存在", "")
		}

		//判断登录信息是否正确
		if util.Md5(password)!=strings.Trim(user.Password," "){
			c.History("密码错误","")
		}

		//如果账号密码无误则更新admin信息
		user.LastIp= c.getClientIp()
		user.LoginCount=user.LoginCount+1
		if _,err:= c.o.Update(&user);err!=nil{
			c.History("登录异常","")
		}else{
			c.History("登录成功","/admin/main.html")
		}
		c.SetSession("user",user)
	}

	//注册登录模板
	c.TplName= c.controllerName+"/login.html"

}

//退出登录
func (c *AdminController) Logout() {
	c.DestroySession()
	c.History("退出登录","/admin/login.html")
}

//配置后台
func (c *AdminController) Config(){
	//读取后台信息，后台信息由ID,name,value三列组成
	//key-value对可以是QQ, URL, 时区等内容
	var result []*models.Config
	c.o.QueryTable(new(models.Config).TableName()).All(&result)

	//建两个map
	//option用于向模板中传递参数，和数据库中的表示方法一样：name-value，所以是string-string
	//mp则用与检查，表示方法是：name-Config{}结构体
	options := make(map[string]string)
	mp := make(map[string]*models.Config)
	for _, v := range result {
		options[v.Name] = v.Value
		mp[v.Name] = v //v是一个Config结构体
	}
	if c.Ctx.Request.Method == "POST" {
		keys := []string{"url", "title", "keywords", "description", "email", "start", "qq"}
		for _, key := range keys {
			val := c.GetString(key) //从我们post的请求中读取参数
			//先检查一下数据库中有没有我们读取到的参数
			if _, ok := mp[key]; !ok {  //没有
				c.o.Insert(&models.Config{Name: key, Value: val})
			} else { //如果已经有了，则需要更新
				opt := mp[key]
				if _, err := c.o.Update(&models.Config{Id: opt.Id, Name: opt.Name, Value: val/*此处使用POST数据*/}); err != nil {
					panic(err)
				}
			}
		}
		c.History("设置数据成功", "")
	}
	c.Data["config"] = options
	c.TplName = c.controllerName + "/config.html"
}

//后台管理的main页面
func (c *AdminController) Main() { //不要忘了注册后台主页面
	c.TplName = c.controllerName + "/main.tpl"
}

/*************Category*******************/

//后台管理的Category页面
func (c *AdminController) Category(){
	//categories := []*models.Category{} 不要这么声明切片
	//var categories [] *models.Category 这样也不好
	categories:=make([] *models.Category,0)
	c.o.QueryTable(new(models.Category).TableName()).All(&categories)
	c.Data["categorys"]= categories //虽然categorys拼写错误，但是最好不要修改，因为模板里面用的这个名字。:)
	c.TplName=c.controllerName+"/category.tpl"
}

//添加category
func (c *AdminController) Categoryadd(){
	//一般来说可以采用c.GetInt等方法，但如果不是Int64则需要采用下面的办法
	id:=c.GetString("id")
	//id:=c.Input().Get("id")
	if id!=""{ //id为空表示是新添加，不为空则表示修改
		intId,_:=strconv.Atoi(id)
		cate:=models.Category{Id:intId}
		c.o.Read(&cate)
		c.Data["cate"]=cate //引导
	}
	c.TplName = c.controllerName + "/category_add.tpl"
}

//类别的插入和更新
func (c *AdminController) CategorySave(){
	name := c.Input().Get("name") //c.GetString也可以
	id := c.Input().Get("id")
	category:=models.Category{}
	category.Name=name
	if id==""{ //插入
		if _,err:=c.o.Insert(&category);err!=nil{
			c.History("插入数据错误","")
		}else{
			c.History("插入数据成功","/admin/category.html")
		}
	}else{ //更新
		intId,err:=strconv.Atoi(id)
		if err!=nil{
			c.History("参数错误","")
		}
		category.Id=intId
		if _,err:=c.o.Update(&category);err!=nil{
			c.History("更新出错","")
		}else{
			c.History("插入数据成功","/admin/category.html")
		}
	}
}

//类别的删除
func (c *AdminController) CategoryDel(){
	strId:=c.Input().Get("id")
	intId,err:=strconv.Atoi(strId)
	if err!=nil{
		c.History("参数错误","")
	}else{
		if _,err:=c.o.Delete(&models.Category{Id:intId});err!=nil{
			c.History("未能成功删除", "")
		}else {
			c.History("删除成功", "/admin/category.html")
		}
	}
}

/*************Category*******************/

/////////////////////////////////////////
/************Post*********************/
//博客列表
func (c *AdminController) Index(){
	categories:=make([]*models.Category,0)
	c.o.QueryTable(new(models.Category).TableName()).All(&categories)
	c.Data["categorys"]=categories

	var(
		page int
		pagesize int=8
		offset int
		list []*models.Post
		keyword string
		cateId int
	)
	keyword = c.GetString("title")
	cateId,_=c.GetInt("cate_id")
	//第一次访问时，page未初始化，为0
	if page,_=c.GetInt("page");page<1{
		page=1
	}
	offset=(page-1)*pagesize
	q:=c.o.QueryTable(new(models.Post).TableName())
	if keyword!=""{
		q=q.Filter("title__contains",keyword)//包含关键词的title
	}
	count,_:=q.Count()
	if count>0{
		//第一个参数是指定获取几条数据，第二个参数指定从哪里获取start
		//如果页码是1，则访问8篇，从0开始
		//如果页码是2，则访问8篇，从9开始
		q.OrderBy("-created").Limit(pagesize,offset).All(&list)
	}
	//查询过滤数据完毕，向模板写入
	c.Data["keyword"]=keyword
	c.Data["count"]=count
	c.Data["list"]=list
	c.Data["cate_id"]=cateId
	path:=fmt.Sprintf("/admin/index.html?keyword=%s", keyword)
	c.Data["pagebar"]=util.NewPager(page,int(count),pagesize, path,true).ToString()
	c.TplName = c.controllerName + "/list.tpl"
}

//删除博文
func (c *AdminController) Delete() {
	id, err := strconv.Atoi(c.Input().Get("id"));
	if err != nil {
		c.History("参数错误", "")
	}else{
		if _,err := c.o.Delete(&models.Post{Id:id}); err !=nil{
			c.History("未能成功删除", "")
		}else {
			c.History("删除成功", "/admin/index.html")
		}
	}
}

//添加博文
func (c *AdminController) Article() {
	categorys := make([]*models.Category,0)
	c.o.QueryTable( new(models.Category).TableName()).All(&categorys)
	id, _ := c.GetInt("id")
	if id != 0{
		post := models.Post{Id:id}
		c.o.Read(&post)
		c.Data["post"] = post
	}
	c.Data["categorys"] = categorys
	c.TplName = c.controllerName + "/_form.tpl"
}

//保存修改
func (c * AdminController) Save()  {
	post := models.Post{}
	post.UserId = 1
	post.Title = c.Input().Get("title")
	post.Content = c.Input().Get("content")
	post.IsTop,_ = c.GetInt8("is_top")
	post.Types,_ = c.GetInt8("types")
	post.Tags = c.Input().Get("tags")
	post.Url = c.Input().Get("url")
	post.CategoryId, _ = c.GetInt("cate_id")
	post.Info = c.Input().Get("info")
	post.Image = c.Input().Get("image")
	post.Created = time.Now()
	post.Updated = time.Now()

	id ,_ := c.GetInt("id")
	if id == 0 {
		if _, err := c.o.Insert(&post); err != nil {
			c.History("插入数据错误"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/index.html")
		}
	}else {
		post.Id = id
		if _, err := c.o.Update(&post); err != nil {
			c.History("更新数据出错"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/index.html")
		}
	}
}

/************Post*********************/
