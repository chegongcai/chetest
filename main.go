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
		if buf[n-1] != '#' {
			return
		}
		rev_buf := string(buf[0 : n-2]) //delete the tail #
		ParseProtocol(rev_buf, conn)    //do protocol parse
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

func DeleteTail(buf string) string {
	lenght := len([]rune(buf))
	data := []byte(buf)
	str := string(data[0 : lenght-1])
	return str
}

func testbuf() {

	var temp []string
	var flag string = "hello,che,123n,uio"

	temp = strings.Split(flag, ",")
	buf := DeleteTail(string(temp[2]))
	fmt.Println(buf)
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf []string

	fmt.Println("Receive from client", rev_buf)

	arr_buf = strings.Split(rev_buf, ",")

	fmt.Println(arr_buf[0])

	switch arr_buf[0] {
	case "BDT01":
		fmt.Println("IMEI:", arr_buf[1])
		zone, _ := strconv.Atoi(GetZone())
		buf := fmt.Sprintf("BDS01,%s,%d#", GetTimeStamp(), zone)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT02":
		fmt.Println("time stamp:%s%s", string(arr_buf[1]), string(arr_buf[5]))
		latitude := DeleteTail(string(arr_buf[2]))
		longtitude := DeleteTail(string(arr_buf[3]))
		fmt.Println("latitude:", latitude)
		fmt.Println("longtitude:", longtitude)
		fmt.Println("speed:", arr_buf[4])
		fmt.Println("angle:", arr_buf[6])
		fmt.Println("data:", arr_buf[7])
		buf := fmt.Sprintf("BDS02#")
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT03":

		break
	}
	if err != nil {
		return
	}
}
