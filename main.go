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
	var rev_buf *string
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		//fmt.Println("Receive from client", rAddr.String(), string(buf[0:n]))
		fmt.Println("client IP", rAddr.String())
		*rev_buf = string(buf[0:n])
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

func GetAsciiStrFromBuffer(StrOut1 *string, buf *string, strMaxlen int, StrIn *string) int {
	i := 0

	str := []byte(*StrIn)
	var strBuf [512]byte

	fmt.Println(string(str[0:]))

	for i = 0; (str[i] != ',') && (str[i] != '#'); i++ {
		strBuf[i] = str[i]

		if i >= strMaxlen {
			return 0
		}
	}
	i++
	*StrOut1 = string(strBuf[0:i])
	*buf = string(str[i:])
	return i
}

func testbuf() {
	var command, test, buf, imei, zone *string

	test = new(string)
	command = new(string)
	buf = new(string)
	imei = new(string)
	zone = new(string)

	*test = "BDS01,20180718143213,8#"

	GetAsciiStrFromBuffer(command, buf, 6, test)
	fmt.Println("get command:", *command)

	GetAsciiStrFromBuffer(imei, buf, 15, buf)
	fmt.Println("get imei:", *imei)

	GetAsciiStrFromBuffer(zone, buf, 1, buf)
	fmt.Println("get zone:", *zone)
}

func ParseProtocol(rev_buf *string, conn net.Conn) {
	var err error
	var command, buf *string
	fmt.Println("Receive from client", rev_buf)

	GetAsciiStrFromBuffer(command, buf, 6, rev_buf)
	fmt.Println("get command:", *command)

	switch *command {
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
