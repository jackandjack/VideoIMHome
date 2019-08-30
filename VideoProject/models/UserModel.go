package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/dgrijalva/jwt-go/request"
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
	//输出数据库日记
	orm.Debug = true
	fmt.Sprintln("db", db)
}

/*添加用户接口*/
func AddUser(user User) bool {

	id, err := orm.NewOrm().Insert(&user)
	fmt.Println("id", id)
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

	db, _ := orm.NewQueryBuilder("mysql")
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

func Regist(user User, username string, password string) {

}

/*登录接口*/
func Login(user User, number string, password string) bool {

	db, _ := orm.NewQueryBuilder("mysql")
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
func Queryall(user *User) (*User, bool) {

	var users []*User
	db, _ := orm.NewQueryBuilder("mysql")
	db.Select("*").From("user").Where(fmt.Sprintf("username='%s'", user.Username)).And(fmt.Sprintf("Password='%s'", user.Password))
	mysql := db.String()
	userCount, err := orm.NewOrm().Raw(mysql).QueryRows(&users)
	fmt.Println("count=", userCount)
	for _, u := range users {
		if user.Username == u.Username {
			user = u
			fmt.Println("username=", user.Username)
		}
	}
	if (err == nil) && (int(userCount) > 0) {
		return user, true
	} else {
		return user, false
	}

}

type User struct {
	Id           int
	Nickname     string `orm:"size(100)""`
	Username     string `orm:"size(100)"`
	Password     string `orm:"size(100)"`
	Number       string `orm:"size(100)"`
	Introduction string `orm:"size(100)"`
	SchoolNum    string `orm:"size(100)"`
	Sex          int64  `orm:"default(1)"`
	Birthday     string `orm:"size(100)"`
	Region       string `orm:"size(100)"`
	Token        string `orm:"token"`
	Praise       int    `orm:"default(1)"`
	Focus        int    `orm:"default(1)"`
	fans         int    `orm:"default(1)"`
	Icon         string `orm:"default(http://hk-pipixie.oss-cn-hongkong.aliyuncs.com/1564975612_4.png)"`
}
