package routers

import (
	"HelloBee/controller"

	"github.com/astaxie/beego/context" //此处非过滤器函数默认导入函数(默认导入函数错误)

	"github.com/astaxie/beego"
)

func init() {
	//路由过滤器
	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc) //使用过滤器函数必须在必须在配置文件中开启sessionon = true 写入
	//注册路由
	beego.Router("/register", &controller.UserController{}, "get:ShowRegister;post:HandlePost")
	//登录路由
	beego.Router("/login", &controller.UserController{}, "get:ShowLogin;post:HandleLogin")
	//文章列表页访问
	beego.Router("/Article/showArticleList", &controller.ArticleController{}, "get:ShowArticleList;post:HandleSelect")
	//添加文章
	beego.Router("/Article/addArticle", &controller.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	//显示文章详情
	beego.Router("/Article/showArticleDetail", &controller.ArticleController{}, "get:ShowArticleDetail")
	//删除文章
	beego.Router("/Article/DeleteArticle", &controller.ArticleController{}, "get:HandleDelete")
	//修改文章
	beego.Router("/Article/UpDateArticle", &controller.ArticleController{}, "get:ShowUpDataArticle;post:HandleUpDataArticle")
	//添加类型
	beego.Router("/Article/AddArticleType", &controller.ArticleController{}, "get:ShowAddType;post:HandleAddType")
	//退出登录
	beego.Router("/Article/Logout", &controller.ArticleController{}, "get:Logout")
}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
	}
}
