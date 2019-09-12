package models

import (
	_ "github.com/dgrijalva/jwt-go/request"
	_ "github.com/go-sql-driver/mysql" //import your used driver
)

//func init() {
//	//在这里面注册mysql 数据库
//	orm.RegisterDataBase("default", "mysql", "root:wsz888148@tcp(127.0.0.1:3306)/user_db?charset=utf8&loc=Local", 30)
//	// register model
//	orm.RegisterModel(new(Message))
//	// create table
//	orm.RunSyncdb("default", false, true)
//	//初始化数据库
//	//输出数据库日记
//	orm.Debug = true
//}

//type Message struct {
//	MessageId     int64  `orm:"pk;auto;column(id)"`
//	SendToken     string `orm:"size(100)"`  //发表用户的Token
//	Message       string `orm:"size(100)"`  //消息
//	AcceptToken   string `orm:"size(100)"`  //接受用户的Token
//	Type          int                       //1、发送消息 对方未接受 2、对方已经接受
//	SendTime      time.Time                 //发送时间
//}
