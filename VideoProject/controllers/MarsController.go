package controllers

import (
	"VideoIMHome/VideoProject/models"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/proto"
	"net"
)

func init() {

	//listener, err := net.Listen("tcp", "127.0.0.1:8082")
	//
	//if err != nil {
	//
	//	fmt.Println("connet is falier")
	//
	//	return
	//
	//}
	//defer listener.Close()
	//
	//for {
	//	conne, err := listener.Accept()
	//	if err != nil {
	//		fmt.Println("err = ", err)
	//	}
	//    go HandleConect(conne)
	//}
}

type Msg struct {
	HeadLength    uint32
	ClientVersion uint32
	CmdId         uint32
	Seq           uint32
	BodyLength    uint32
}

func HandleConect(conn net.Conn) {
	//函数调用关闭conne
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println("sucess=", addr)
	//读取用户数据
	var buf = make([]byte, 2048)
	isQuit := make(chan bool)

	n, err := conn.Read(buf)

	if n == 0 {

		isQuit <- true

		fmt.Println("客户端退出Socket")
	}

	if err != nil {
		fmt.Println("读取数据是失败=", err)
		return
	}
	msgH := Msg{}
	newbuf := bytes.NewBuffer([]byte(string(buf[:20])))
	binary.Read(newbuf, binary.BigEndian, &msgH)
	if msgH.CmdId == 6 {
		fmt.Println("该信息是心跳body=", msgH.BodyLength)
		conn.Write([]byte(string(buf[:n])))
	} else {

		//获取Body 的Buffer 然后转Protobuf

		Bodybuf := bytes.NewBuffer([]byte(string(buf[20:n])))
		var bodyPro = &models.HelloRequest{}
		proto.Unmarshal(Bodybuf.Bytes(), bodyPro)
		fmt.Println("bodyPro.user=", bodyPro.User)
		fmt.Println("bodyPro.Text=", bodyPro.Text)
		var buffer bytes.Buffer
		//获取得到 headbuf 的数据显示
		resBody := &models.HelloResponse{
			Retcode: 1000,
			Errmsg:  "消息发成功",
		}
		data, _ := proto.Marshal(resBody)
		bufHead := &bytes.Buffer{}
		err := binary.Write(bufHead, binary.BigEndian, msgH)
		if err != nil {
			panic(err)
		}
		buffer.Write(bufHead.Bytes())
		buffer.Write(data)
		fmt.Println("buffer=", buffer.String())
		_, errmsg := conn.Write(buffer.Bytes())
		if errmsg != nil {
			beego.Error("send err=", errmsg)
		}

	}
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
