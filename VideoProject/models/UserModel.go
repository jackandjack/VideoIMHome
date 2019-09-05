package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go/request"
	_ "github.com/go-sql-driver/mysql" //import your used driver
	"time"
)

const (
	SecretKey = "VideoProgress_2016"
)

var (
	db orm.QueryBuilder
)

func init() {
	//在这里面注册mysql 数据库
	orm.RegisterDataBase("default", "mysql", "root:wsz888148@tcp(127.0.0.1:3306)/user_db?charset=utf8&loc=Local", 30)
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
func AddUser(user User) (int64, bool) {
	id, err := orm.NewOrm().Insert(&user)
	fmt.Println("id", id)
	if err == nil {
		return id, true
	} else {
		return 0, false
	}
}

/*修改个人资料*/
func UpdateUser(user User, key string, value string, phone string) bool {
	db, _ := orm.NewQueryBuilder("mysql")
	db.Update("user").Set(fmt.Sprintf("%s='%s'", key, value)).Where(fmt.Sprintf("phone='%s'", phone))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err == nil {
		return false
	} else {
		return true
	}
}
func GetToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(SecretKey))
	fmt.Println("Token=", tokenString)
	return tokenString
}
func Regist(user User) (User, bool) {
	t := time.Now()
	fmt.Println(t)
	t1 := t.Unix()
	var times = fmt.Sprintf("%d", t1)
	user.Token = GetToken()
	user.Lasttime = times
	user.Logintime = times
	user.Introduction = "这个人很懒,什么都没有写"
	user.Nickname = fmt.Sprintf("项目_%s", times)
	user.Icon = "http://hk-pipixie.oss-cn-hongkong.aliyuncs.com/1564975612_4.png"
	//先查找phone 是不是存在
	newUser, iscuess := Findphone(user.Phone)
	if iscuess {
		return newUser, true
	} else {
		id, sucess := AddUser(user)
		if sucess {
			return QueryUser(id), true
		} else {
			return QueryUser(id), false
		}
	}
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
func Findphone(phone string) (User, bool) {
	var user User
	db, _ := orm.NewQueryBuilder("mysql")
	db.Select("*").From("user").Where(fmt.Sprintf("phone=%s", phone))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err != nil {
		return user, false
	} else {
		return user, true
	}
}
func QueryUser(id int64) User {
	var user User
	db, _ := orm.NewQueryBuilder("mysql")
	db.Select("*").From("user").Where(fmt.Sprintf("id=%d", id))
	mysql := db.String()
	err := orm.NewOrm().Raw(mysql).QueryRow(&user)
	if err != nil {
		return user
	} else {
		return user
	}
}

type User struct {
	Id           int
	Nickname     string `orm:"size(100)""` //昵称
	Phone        string `orm:"size(100)"`  //电话
	Introduction string `orm:"size(100)"`  //说明--简介
	Sex          int64  `orm:"default(0)"` //性别 0 保密 1 男 2 女
	Birthday     string `orm:"size(100)"`  //生日 1993 -07 -11
	Region       string `orm:"size(100)"`  //地区 注册当前日期
	Token        string `orm:"token"`      //用户的Token
	Praise       int    `orm:"default(1)"` //点赞数
	Focus        int    `orm:"default(1)"` //关注数
	fans         int    `orm:"default(1)"` //粉丝数
	Icon         string //头像图片
	Lasttime     string //时间戳 当前的注册时间
	Logintime    string //时间戳 当前的登录时间
	Type         string // 注册的设备 1、iOS 2、安卓
}
