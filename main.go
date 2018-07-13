/*
package main

import (
	_ "chetest/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.Info("beego.Run!!!!")
	beego.Run()
}
*/
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		str_command := string(buf[38:40])
		if buf[0] != '[' && buf[n-1] != ']' {
			fmt.Println("data error!!!!!!!!!!")
			return
		}
		ParseProtocol(str_command, conn) //do protocol parse
	}
}

func GetTimeStamp() string {
	buf := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	return buf
}

func ParseProtocol(command string, conn net.Conn) {
	var err error
	switch command {
	case "T0":
		buf := fmt.Sprintf("%s,S0", GetTimeStamp())
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "T1":
		buf := fmt.Sprintf("%s,S1,1", GetTimeStamp())
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "T3":
		buf := fmt.Sprintf("%s,S3", GetTimeStamp())
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	}
	if err != nil {
		return
	}
}
