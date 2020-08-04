package model

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//表的设计

//定义一个结构体
type User struct {
	Id       int
	Name     string
	PassWord string
	//Pass_Word
	Articles []*Article `orm:"rel(m2m)"` //多对多（n）
}

//文章表和文章类型表：一对多
type Article struct {
	Id          int          `orm:"pk;auto"`
	ArtiName    string       `orm:"size(20)"`
	Atime       time.Time    `orm:"auto_now"`
	Acount      int          `orm:"default(0);null"`
	Acontent    string       `orm:"size(500)"`
	Aimg        string       `orm:"size(100)"`
	ArticleType *ArticleType `orm:"rel(fk)"`       //一对多（1）
	Users       []*User      `orm:"reverse(many)"` //多对多（n）
}

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` //一对多（n）
}

//在ORM里面__是有特殊含义的

func init() {
	//ORM操作数据库
	//获取连接对象
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")

	//创建表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	//生成表
	//第一个参数是数据库别名，第二个参数是是否强制更新
	orm.RunSyncdb("default", false, true)
	//操作表

}
