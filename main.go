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
		ParseProtocol(str_command, conn) //do protocol parse
	}
}

func ParseProtocol(command string, conn net.Conn) {
	switch command {
	case "T1":
		buf := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d,S1,1", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		fmt.Println(buf)
		_, err2 := conn.Write([]byte(buf))
		if err2 != nil {
			return
		}
		break
	case "T3":
		buf := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d,S3", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		fmt.Println(buf)
		_, err2 := conn.Write([]byte(buf))
		if err2 != nil {
			return
		}
		break
	}
}
