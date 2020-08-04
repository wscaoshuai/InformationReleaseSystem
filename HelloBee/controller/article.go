package controller

import (
	"HelloBee/model"
	"math"
	"path"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页(文章分页、下拉选择框)
func (this *ArticleController) ShowArticleList() {
	//获取数据
	//1.查询
	typeName := this.GetString("select")
	//创建一个orm对象
	o := orm.NewOrm()
	qs := o.QueryTable("Article") //QueryTable(//高级查询:指定表)的返回值：queryseter
	var articles []model.Article
	// _, err := qs.All(&articles) //select *from articles
	// if err != nil {
	// 	beego.Info("查询数据错误")
	// }
	// beego.Info(articles[0])

	//查询总记录数
	count, _ := qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	//获取总页数
	pageSize := 2
	pageCount := math.Ceil(float64(count) / float64(pageSize)) //向上取整函数：math.Ceil
	//获取页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	//获取数据
	//作用就是获取数据库部分数据,第一个参数，获取几条,第二个参数，从那条数据开始获取,返回值还是querySeter
	//起始位置计算
	start := (pageIndex - 1) * pageSize
	//1.pagesize一页显示多少；2.start起始位置；3.all(存放位置)
	//RelatedSel关联查询ArticleType,如果没有数据就不显示(则需要在查询总数目时，同时查询ArticleType)
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	//获取类型数据
	var types []model.ArticleType
	_, err = o.QueryTable("ArticleType").All(&types)
	if err != nil {
		beego.Info("获取数据错误")
	}

	//根据类型获取数据
	//1.接收数据
	//typeName := this.GetString("select")
	//beego.Info(TypeName)
	//2.处理数据
	var articleswithtype []model.Article
	if typeName == "" {
		beego.Info("下拉框传递数据失败")
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articleswithtype)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articleswithtype)
	}
	//3.查询数据
	//此处需要Filter("查询字段","要匹配的值")过滤器函数相当于：select * form user where↓
	//因为Filter是惰性查询，所以需要RelatedSel指定关联的查询表格

	userName := this.GetSession("userName")

	//传递数据
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articleswithtype
	this.Data["types"] = types
	this.Data["typeName"] = typeName
	this.Data["userName"] = userName
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//处理下拉框改变请求
func (this *ArticleController) HandleSelect() {
	//1.接收数据
	TypeName := this.GetString("select")
	//beego.Info(TypeName)
	//2.处理数据
	if TypeName == "" {
		beego.Info("下拉框传递数据失败")
		return
	}
	//3.查询数据
	o := orm.NewOrm()
	var articles []model.Article
	//此处需要Filter("查询字段","要匹配的值")过滤器函数相当于：select * form user where↓
	//因为Filter是惰性查询，所以需要RelatedSel指定关联的查询表格
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", TypeName).All(&articles)
	//beego.Info(articles)
}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//查询类型数据，传递到视图中
	o := orm.NewOrm()
	var types []model.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "add.html"

}

//获取添加文章数据
func (this *ArticleController) HandleAddArticle() {
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	//2校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}

	//处理文件上传
	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}

	//1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return
	}

	//2.文件格式
	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return
	}

	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储
	this.SaveToFile("uploadname", "../static/img/.."+fileName)
	if err != nil {
		beego.Info("上传文件失败")
		return
	}

	//3.处理数据
	//3.1创建连接
	o := orm.NewOrm()

	article := model.Article{}
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "../static/img/.." + fileName

	//下拉框
	//给Article对象赋值
	//获取到下拉框传递过来的类型数据
	typeName := this.GetString("select")
	//类型判断
	if typeName == "" {
		beego.Info("下拉框数据错误")
		return
	}
	//获取type对象
	var artiType model.ArticleType
	artiType.TypeName = typeName
	err = o.Read(&artiType, "TypeName")
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	article.ArticleType = &artiType

	//3.2插入操作
	_, err = o.Insert(&article)
	if err != nil {
		beego.Info("插入数据失败")
		return
	}

	//4.返回页面
	this.Redirect("/Article/showArticleList", 302)
}

//展示文章详情页面
func (this *ArticleController) ShowArticleDetail() {
	//获取数据
	id, err := this.GetInt("articleId")
	//数据校验
	if err != nil {
		beego.Info("传递的链接错误")
	}
	//操作数据
	o := orm.NewOrm()
	// var article model.Article
	// article.Id = id
	article := model.Article{Id: id}
	err = o.Read(&article)
	if err != nil {
		beego.Info("查询错误", err)
		return
	}

	//修改阅读量
	article.Acount += 1

	//多对多插入读者
	//1.获取操作对象
	// artile := model.Article{Id:id}//前面已经获取过了，所以无需重复
	//2.获取多对多操作对象
	m2m := o.QueryM2M(&article, "Users")
	//3.获取插入对象
	userName := this.GetSession("userName")
	user := model.User{}
	user.Name = userName.(string) //此处Session获取的是interface，所以此处要进行断言
	o.Read(&user, "Name")
	//4.多对多插入
	_, err = m2m.Add(&user)
	if err != nil {
		beego.Info("插入失败")
		return
	}

	//多对多的第一种方式
	//o.LoadRelated(&article, "Users")

	//多对多的第二种方式
	var users []model.User
	beego.Info(article)
	//错误:o.QueryTable("Article").Filter("Users__User__UserName",userName.(string)).Distinct().Filter("Id",id).One(&article)
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)

	//返回视图页面
	o.Update(&article) //更新视图，没有指定的话，会自己进行查询
	this.Data["users"] = users
	this.Data["article"] = article
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["detailHead"] = "head.html"
	this.TplName = "content.html"
}

//删除文章
func (this *ArticleController) HandleDelete() {
	//1.URL传值
	id, _ := this.GetInt("id")
	//1.1orm对象
	o := orm.NewOrm()
	//1.2删除对象
	article := model.Article{Id: id}
	//1.3删除
	o.Delete(&article)
	this.Redirect("/Article/showArticleList", 302)

}

//展示修改文章页面
func (this *ArticleController) ShowUpDataArticle() {
	//获取数据
	//1.URL传值
	id := this.GetString("id")
	//2.判断数据
	if id == "" {
		beego.Info("连接错误")
		return
	}
	//3.查询操作
	o := orm.NewOrm()
	article := model.Article{}
	//类型转换
	id2, _ := strconv.Atoi(id)
	article.Id = id2
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询错误")
		return
	}
	//把数据传递给视图
	this.Data["article"] = article
	this.TplName = "updata.html"
}

//修改文章
func (this *ArticleController) HandleUpDataArticle() {
	//1.拿数据
	name := this.GetString("articleName")
	content := this.GetString("content")
	id, _ := this.GetInt("id")
	//2.判断数据
	if name == "" || content == "" {
		beego.Info("更新数据失败")
		return
	}
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		beego.Info("上传文件失败")
		return
	}
	defer file.Close() //关闭文件！！！

	//2.1判断大小
	if head.Size > 50000 {
		beego.Info("图片大小错误")
		return
	}
	//2.2判断类型
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != "png" && ext != "jpeg" {
		beego.Info("图片格式错误")
		return
	}
	//2.3防止文件名重复
	filename := time.Now().Format("2006-01-02-15:03:04")
	this.SaveToFile("uploadname", "./static/img/"+filename+ext) //（"文件上传文件名"，"路径和文件名"，""）

	//3.更新操作
	//3.1获取orm对象
	o := orm.NewOrm()
	//3.2创建更新对象
	article := model.Article{Id: id}
	//3.3读取操作
	err = o.Read(&article)
	if err != nil {
		beego.Info("要更新的文章不存在")
		return
	}
	//3.4更新
	article.ArtiName = name
	article.Acontent = content
	article.Aimg = "../static/img/" + filename + ext
	_, err = o.Update(&article)
	if err != nil {
		beego.Info("更新失败")
		return
	}

	//4.跳转页面
	this.Redirect("/Article/showArticleList", 302)
}

//展示添加文章类型
func (this *ArticleController) ShowAddType() {
	//1.读取类型表，显示数据
	o := orm.NewOrm()
	var ArticleTypes []model.ArticleType
	_, err := o.QueryTable("ArticleType").All(&ArticleTypes)
	if err != nil {
		beego.Info("查询类型错误")
	}
	this.Data["types"] = ArticleTypes

	this.TplName = "addType.html"
}

//文章类型
func (this *ArticleController) HandleAddType() {
	//1.拿数据
	typename := this.GetString("typeName")
	//2.判断数据
	if typename == "" {
		beego.Info("添加数据类型为空")
		return
	}
	//3.执行插入操作
	o := orm.NewOrm()
	var artiType model.ArticleType
	artiType.TypeName = typename
	_, err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	// //4.执行删除操作
	// //4.1URL传值
	// typedel := this.GetString("delName")
	// if typedel == "" {
	// 	beego.Info("获取数据失败")
	// }
	// //4.2删除对象
	// article := model.Article{Id: typedel}
	// //4.3删除
	// err = o.Delete(&article)
	// if err != nil {
	// 	beego.Info("删除失败")
	// 	return
	// }

	//5.展示视图
	this.Redirect("/Article/AddArticleType", 302)
}

//退出登录
func (this *ArticleController) Logout() {
	//1.退出登录状态(删除session)
	this.DelSession("userName")
	//2.跳转登录界面
	this.Redirect("/login", 302)
}
