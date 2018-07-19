package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//182.254.185.142  8080

func main() {
	service := ":8080"
	//testbuf()
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
		fmt.Println("client IP", rAddr.String())
		rev_buf := string(buf[0:n])
		ParseProtocol(rev_buf, conn) //do protocol parse
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

func testbuf() {

	var temp []string
	var flag string = "hello,che,123,uio"

	temp = strings.Split(flag, ",")

	fmt.Println(temp[3])
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf []string

	fmt.Println("Receive from client", rev_buf)

	arr_buf = strings.Split(rev_buf, ",")

	fmt.Println(arr_buf[0])

	switch arr_buf[0] {
	case "BDT01":
		zone, _ := strconv.Atoi(GetZone())
		buf := fmt.Sprintf("BDS01,%s,%d#", GetTimeStamp(), zone)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT02":

		break
	case "BDT03":

		break
	}
	if err != nil {
		return
	}
}

/*
func ParseProtocol(rev_buf *string, conn net.Conn) {
	var err error
	var command, buf_res *string
	var bdy int

	command = new(string)
	buf_res = new(string)

	fmt.Println("Receive from client", *rev_buf)

	GetAsciiStrFromBuffer(command, buf_res, 6, rev_buf)
	fmt.Println("get command:", *command)

	if strings.Contains(*command, "BDT01") == true {
		bdy = 1
		fmt.Println("bdy = 1")
	}
	switch bdy {
	case 1:
		var imei *string
		imei = new(string)
		GetAsciiStrFromBuffer(imei, buf_res, 15, buf_res)
		fmt.Println("get imei:", *imei)
		zone, _ := strconv.Atoi(GetZone())
		buf := fmt.Sprintf("BDS01,%s,%d#", GetTimeStamp(), zone)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case 2:

		break
	case 3:

		break
	}
	if err != nil {
		return
	}
}
*/
