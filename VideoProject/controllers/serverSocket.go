package controllers

/**
* 这是个简易的游戏服务器
* 每条包结构会有个包头标志剩下的包体长度
* 会进行断包粘包处理.
* 每个客户端由一个Client实例来维护
 */

import (
	"VideoIMHome/VideoProject/models"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

var (
	ClientMap map[int]net.Conn = make(map[int]net.Conn)

	Clients map[int]Client = make(map[int]Client)
)

func init() {
	//fmt.Println("CPU个数: ", runtime.NumCPU())
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//fmt.Println(os.Args[0])
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8083")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	clientIndex := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("监听错误: ", err.Error())
			continue
		}
		clientIndex++
		go handleClient(conn, clientIndex)
	}
}

type Client struct {
	conn net.Conn
}

func (this *Client) quit() {
	this.conn.Close()
}

func handleClient(conn net.Conn, index int) {
	c := Client{conn}
	Clients[index] = c
	ClientMap[index] = conn
	fmt.Println("新用户连接, 来自: ", conn.RemoteAddr(), "index: ", index)
	isHeadLoaded := false
	reader := bufio.NewReader(conn)
	defer func() {
		conn.Close()
		delete(ClientMap, index)
		fmt.Println("移除序号为: ", index, "的客户端")
	}()

Out:
	for {
		if !isHeadLoaded {
			headLenSl := make([]byte, 1024*1024)
			//已经读取的包头字节数
			Buflen, err := reader.Read(headLenSl)
			if err != nil {
				fmt.Println("客户端断开连接", err)
				break Out
			}
			if Buflen < 20 {
				break Out
			}
			msgH := HeadMsg{}
			//取前20 获取Head 数据包
			newbuf := bytes.NewBuffer([]byte(string(headLenSl[0:20])))
			binary.Read(newbuf, binary.BigEndian, &msgH)

			if msgH.CmdId != 6 {
				isHeadLoaded = true
			} else {
				sendMsgToAll(headLenSl[:20])
			}
			if isHeadLoaded {
				if Buflen != int(20+msgH.BodyLength) {
					fmt.Println("客户端发送的结构体错误")
					break Out
				}
				fmt.Println("接受到用户数据和消息")
				var bodyPro = &models.HelloRequest{}
				proto.Unmarshal([]byte(string(headLenSl[20:Buflen])), bodyPro)
				fmt.Println("send user=", bodyPro.User)
				resBody := &models.HelloResponse{
					Retcode: proto.Int32(1000),
					Errmsg:  proto.String("发送成功"),
				}
				data, _ := proto.Marshal(resBody)
				msgH.BodyLength = uint32(len(data))
				var sendBufer = bytes.NewBuffer([]byte{})
				sendBufer.Reset()
				sendBufer.Write(Unpack(&msgH))
				sendBufer.Write(data)
				n, senderr := conn.Write(sendBufer.Bytes())
				if senderr != nil {
					fmt.Println("发送失败")
				} else {
					fmt.Println("发送字节长度n=", n)
				}
				isHeadLoaded = false
			}
		}
	}
}

//整型转换为字节
func IntToBytes(n int) []byte {
	data := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	//将data参数里面包含的数据写入到bytesBuffer中
	binary.Write(bytesBuffer, binary.BigEndian, data)

	return bytesBuffer.Bytes()
}

func Unpack(msg *HeadMsg) []byte {

	databufio := bytes.NewBuffer([]byte{})

	binary.Write(databufio, binary.BigEndian, msg.HeadLength)

	binary.Write(databufio, binary.BigEndian, msg.ClientVersion)

	binary.Write(databufio, binary.BigEndian, msg.CmdId)

	binary.Write(databufio, binary.BigEndian, msg.Seq)

	binary.Write(databufio, binary.BigEndian, msg.BodyLength)

	return databufio.Bytes()
}

//func ReciveMsg(msg []byte) proto.Message {
//
//
//	return proto.Message;
//
//}

func sendMsgToAll(msg []byte) {
	for _, value := range ClientMap {
		writer := bufio.NewWriter(value)
		//写入2个字节字符串长度.以供flash读取便利
		writer.Write(msg)
		//writer.WriteString(msg)
		writer.Flush()
	}
}

type HeadMsg struct {
	HeadLength    uint32
	ClientVersion uint32
	CmdId         uint32
	Seq           uint32
	BodyLength    uint32
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
