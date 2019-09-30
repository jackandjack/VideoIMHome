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
	"runtime"
)

var (
	ClientMap map[int]net.Conn = make(map[int]net.Conn)

	Clients map[int]Client = make(map[int]Client)
)

func init() {
	fmt.Println("CPU个数: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
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
			headLenSl := make([]byte, 2480)
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
			newbuf := bytes.NewBuffer([]byte(string(headLenSl[:20])))
			binary.Read(newbuf, binary.BigEndian, &msgH)
			fmt.Println("cmdid: ", msgH.CmdId)
			sendMsgToAll(headLenSl[:20])
			if msgH.CmdId != 6 {
				isHeadLoaded = true
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
					Retcode: 1000,
					Errmsg:  "发送成功",
				}
				data, _ := proto.Marshal(resBody)
				databuf := bytes.NewBuffer(data)

				newmsgH := &HeadMsg{
					HeadLength:    20,
					ClientVersion: 200,
					CmdId:         msgH.CmdId,
					Seq:           msgH.Seq,
					BodyLength:    uint32(databuf.Len()),
				}
				var bufferhead = make([]byte, 20)
				dateBuffer := bytes.NewBuffer(bufferhead)
				err := binary.Write(dateBuffer, binary.BigEndian, newmsgH)
				if err != nil {
					panic(err)
				}
				conn.Write(BytesCombine(dateBuffer.Bytes(), data))
				isHeadLoaded = false
			}
		}
	}
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
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

func parseData(data []byte) {

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
