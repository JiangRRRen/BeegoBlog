:partly_sunny:来源是github用户 **lightTrace** 提供的开源博客项目。我选取了其中有关于**后台管理**的部分，主要的工作是：

- 功能拆分
- 添加详细的注释
- 撰写引导文档（README）

:partly_sunny:发布这个项目的目的是给初学Beego框架的新手提供一个 **短、平、快** 的练手Demo。因此，我把原来的代码拆分为5大块内容，每一块都实现了一个小功能，方便，希望能将学习曲线变得更加平滑。

:moon:使用注意：

1. 博客是基于Beego编写的，请确保你已经学习过Beego框架，且能够运行。

2. 前端部分放在view, static文件夹内，由于前端渲染不是我们的重点，所以直接复制过来用即可。

3. 使用时注意路径，我是放在Blog文件夹下面，请统一修改import的路径，确保正确。

   <img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207142518235.png" alt="image-20200207142518235" style="zoom:67%;" />

4. 项目的结构是：功能层层叠加，wholesrc文件是完整代码。

# 1. 后台登录

## 1.1 实现效果

运行后在浏览器输入`localhost:8080/admin/login`跳出如下界面：

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207130346.png"  style="zoom:45%;display: block; margin: 0px auto; vertical-align: middle;">

输入账号admin，密码123456，登陆成功后显示：

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207131139.png"  style="zoom:45%;display: block; margin: 0px auto; vertical-align: middle;">

然后我们关闭网页，输入`localhost:8080/admin/main`依然可以访问这个页面，证明session有效。之后关闭浏览器，重新打开，输入`localhost:8080/admin/main`发现跳转到登录界面，说明非法访问拦截生效。

## 1.2 实现过程

登录的信息是从数据库来的，所以我们第一步就是需要顺利连接上数据库。

1. 在conf文件下，将数据库的连接信息输入，同时记得开启session
2. 在model文件下，新建一个go文件，用于mysql初始化。

连接上后，我们还需要保证数据收发的畅通。

1. 在models文件下，新建一个go，构建User结构体，用于承载信息。
2. 在models文件下的mysql初始化文件中，将mysql数据库注册到ORM中，同时将User结构体也注册到ORM中。

之后在controllers文件夹下下面添加一个**admin.go作为admin的控制器**，再添加一个**base.go作为控制器辅助组件。**

admin控制器添加以下方法：

- Login登录
- Logout登出
- Main显示后台页面

同时也需要在routers文加下，**注册一个自动路由**。

# 2. 后台信息设置

## 2.1 实现效果

登录后，在系统设置一栏可以看到：

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207134611.png"  style="zoom:45%;display: block; margin: 0px auto; vertical-align: middle;">

修改内容后，显示：

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207134754.png"  style="zoom:45%;display: block; margin: 0px auto; vertical-align: middle;">

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207134823.png"  style="zoom:55%;display: block; margin: 0px auto; vertical-align: middle;">

## 2.2 实现过程

相较于登录部分，我们已经完成了数据库的连接。

首先需要在models文件下添加一个config，用于设置信息的载体。然后在控制器的base.go中添加一个History函数用于显示消息和跳转URL。

最重要的是在admin.go中实现Config()方法，用于配置后台。

新增内容结构如下：

```
--controllers
	--admin.go
		*Config() //控制器方法
    --base.go
    	*History() //用于发送消息和跳转
--models
	--InitModel.go
		*RegisterModel(new(Config)) //在ORM中注册config结构体
	--md_config.go
		*Config struct //创建config结构体
		*TableName() //获取表名
```

# 3. 类别

## 3.1 实现效果

分类管理界面能查看已添加的文章类别。

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207135614.png"  style="zoom:35%;display: block; margin: 0px auto; vertical-align: middle;">

为这些类别实现增删改的功能

<img src="https://uk-1259555870.cos.eu-frankfurt.myqcloud.com/20200207135712.png"  style="zoom:75%;display: block; margin: 0px auto; vertical-align: middle;">

## 3.2 实现过程

思路和2大同小异，首先需要在models文件下添加一个config，用于设置信息的载体。然后在admin.go中设置四个方法：

- category
- categoryadd
- categorySave
- categoryDel

category函数用于显示类别页面，add用于引导修改页面：

<img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207100406781.png" style="zoom:50%;" />

添加和修改按钮会导引向上面的修改页面，然后点击**提交**后，会通过categorySave函数保存，categorySave则是具体功能的实现。

<img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207100305038.png" alt="image-20200207100305038" style="zoom:50%;" />

删除不需要这么麻烦，直接从数据库中清除即可。

# 4. 博文列表

## 4.1 实现效果

点击博文列表，可以查看博文列表。

<img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207141225752.png" alt="image-20200207141225752" style="zoom:50%;" />

## 4.2 实现过程

除了常规的在models中添加结构体，在admin.go中注册方法外，还需要在util文件下建立一个`newpager.go`的文件，用于页面的操作。

- Pager struct
- url function
- ToString function

# 5. 博文修改

## 5.1 实现效果

这一步是对前面博文列表的补充，实现博文的添加、修改和删除功能。

<img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207141506393.png" alt="image-20200207141506393" style="zoom:50%;" />

## 5.2 实现过程

博文列表的难点在于页面的逻辑，而修改这一部分不涉及这些，和类别的修改一样，非常简单。需要在admin.go控制器中添加如下函数即可：

<img src="C:\Users\44825\AppData\Roaming\Typora\typora-user-images\image-20200207141617615.png" alt="image-20200207141617615" style="zoom: 67%;" />