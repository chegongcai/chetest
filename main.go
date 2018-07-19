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
const version = 0

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

	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		rAddr := conn.RemoteAddr()
		fmt.Println("client IP", rAddr.String(), n, string(buf[0:n]))
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

func DeleteTail(str string) string {
	lenght := len([]rune(str))
	buf := []byte(str)
	str_out := string(buf[0 : lenght-1])
	return str_out
}

func ParseStatusData(str string) (signal string, sat_num string, bat string, mode string) {
	buf := []byte(str)
	signal = string(buf[0:3])
	sat_num = string(buf[3:6])
	bat = string(buf[6:9])
	mode = string(buf[10:12])
	return signal, sat_num, bat, mode
}

func testbuf() {

	var temp []string
	var flag string = "123,060009080002"

	temp = strings.Split(flag, ",")
	fmt.Println(temp[1])

	signal, sat_num, bat, mode := ParseStatusData(string(temp[1]))

	fmt.Println(signal, sat_num, bat, mode)
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf []string

	fmt.Println("Receive from client", rev_buf)
	arr_buf = strings.Split(rev_buf, ",")
	fmt.Println(arr_buf[0])
	switch arr_buf[0] {
	case "BDT01":
		//parse data
		imei := string(arr_buf[1])
		//printf data
		if version == 0 {
			fmt.Println("IMEI:", imei)
		}
		//send data
		zone, _ := strconv.Atoi(GetZone())
		buf := fmt.Sprintf("BDS01,%s,%d#", GetTimeStamp(), zone)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT02":
		//parse data
		latitude := DeleteTail(string(arr_buf[2]))
		longtitude := DeleteTail(string(arr_buf[3]))
		fmt.Println(arr_buf[7])
		signal, sat_num, bat, mode := ParseStatusData(string(arr_buf[7]))
		speed := string(arr_buf[4])
		angle := string(arr_buf[6])
		//printf data
		if version == 0 {
			fmt.Println("latitude:", latitude)
			fmt.Println("longtitude:", longtitude)
			fmt.Println("speed:", speed)
			fmt.Println("angle:", angle)
			fmt.Println(signal, sat_num, bat, mode)
		}
		//send data
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
