package controllers

import (
	"fmt"
	"net"
)

func init() {

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {

		fmt.Println("connet is falier")

		return

	}
	defer listener.Close()

	for true {
		conne, err := listener.Accept()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		buf := make([]byte, 1024)
		n, err := conne.Read(buf)
		if err != nil {
			fmt.Println("err1=", err)
		}
		fmt.Println("buf=", string(buf[:n]))
		defer conne.Close()
	}
}
