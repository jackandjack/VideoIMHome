package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //import your used driver
)

var (
	db orm.QueryBuilder
)

func init() {
	//在这里面注册mysql 数据库
	orm.RegisterDataBase("default", "mysql", "root:wsz888148@tcp(127.0.0.1:3306)/user_db?charset=utf8", 30)
	// register model
	orm.RegisterModel(new(User))
	// create table
	orm.RunSyncdb("default", false, true)
	//初始化数据库
	db, _ := orm.NewQueryBuilder("mysql")
	//输出数据库日记
	orm.Debug = true
	fmt.Sprintln("db", db)
}

/*添加用户接口*/
func AddUser(user User) bool {
	db.InsertInto("user")
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err == nil {
		beego.Error("插入user失败")
		return false
	} else {
		beego.Info("插入User成功")
		return true
	}
}

/*修改个人资料*/
func UpdateUser(user User, newvalue string, oldvalue string, key string) bool {

	db.Update("user").Set(fmt.Sprintf("%s", newvalue)).Where(fmt.Sprintf("%s='%s'", key, oldvalue))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err == nil {
		beego.Error("更新user失败")
		return false
	} else {
		beego.Info("更新User成功")
		return true
	}
}

/*登录接口*/
func Login(user User, number string, password string) bool {
	db.From("user").Where(fmt.Sprintf("Number='%s'", number)).Where(fmt.Sprintf("Password='%s'", password))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err == nil {
		beego.Error("登录失败")
		return false
	} else {
		beego.Info("登录成功")
		return true
	}
}

/*根据Token 获取用户信息*/
func Queryall(user User) (User, bool) {
	db.From("user").Where(fmt.Sprintf("Token='%s'", user.Token))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err == nil {
		beego.Error("登录失败")
		return user, false
	} else {
		beego.Info("登录成功")
		return user, true
	}
}

type User struct {
	Id           int
	Nickname     string
	Username     string
	Password     string
	Number       string
	Introduction string
	SchoolNum    string
	Sex          string
	Birthday     string
	Region       string
	Token        string
	Praise       string
	Focus        string
	fans         string
	Icon         string
}
