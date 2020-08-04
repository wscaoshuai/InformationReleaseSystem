package controller

import (
	"HelloBee/model"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandlePost() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")

	//beego.Info(userName,pwd)
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整，请重新注册"
		beego.Info("注册数据不完整，请重新注册")
		this.TplName = "register.html"
		return
	}
	//3.操作数据
	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	user := model.User{}
	//给插入对象赋值
	user.Name = userName
	user.PassWord = pwd
	//插入
	o.Insert(&user)
	//返回结果

	//4.返回页面
	//this.Ctx.WriteString("注册成功")
	this.Redirect("/login", 302)
	//this.TplName = "login.html"
}

//展示登录页面
func (this *UserController) ShowLogin() {
	name := this.Ctx.GetCookie("userName")
	if name != "" {
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}

	//this.Data["data"] = "heheh"
	this.TplName = "login.html"
}

func (this *UserController) HandleLogin() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整"
		this.TplName = "login.html"
		return
	}

	//3.操作数据
	//1.获取ORM对象
	o := orm.NewOrm()
	user := model.User{}
	user.Name = userName
	user.PassWord = pwd
	err := o.Read(&user, "Name")
	if err != nil {
		// this.Data["errmsg"] = "用户不存在"
		beego.Info("用户名失效")
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		// this.Data["errmsg"] = "密码错误"
		beego.Info("密码失效")
		this.TplName = "login.html"
		return
	}

	//记住用户名
	check := this.GetString("remember")
	beego.Info(check)
	if check == "on" {
		this.Ctx.SetCookie("userName", userName, time.Second*3600)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	this.SetSession("userName", userName)

	//4.返回页面
	//this.Ctx.WriteString("登录成功")
	this.Redirect("/Article/showArticleList", 302)

}
