package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
		//fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		fmt.Println("client IP", rAddr.String())
		rev_buf := string(buf[0:n])
		str_command := string(buf[0:5])
		ParseProtocol(rev_buf, str_command, conn) //do protocol parse
	}
}

func GetTimeStamp() string {
	buf := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	return buf
}

func GetZone() string {
	local, _ := time.LoadLocation("Local")
	local_str := fmt.Sprintf("%s", time.Now().In(local))
	buf := []byte(local_str)
	return string(buf[32:33])
}
func ParseProtocol(rev_buf string, command string, conn net.Conn) {
	var err error
	fmt.Println("Receive from client", rev_buf)
	switch command {
	case "BDT01":
		zone, _ := strconv.Atoi(GetZone())
		buf := fmt.Sprintf("BDS01,%s,%d#", GetTimeStamp(), zone)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "T1":

		break
	case "T3":

		break
	}
	if err != nil {
		return
	}
}
