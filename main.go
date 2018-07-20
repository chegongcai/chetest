package main

import (
	"chetest/BDYString"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//182.254.185.142  8080
const version = 0 // 0 for debug

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
		fmt.Println("client IP", rAddr.String())
		if buf[n-1] != '$' {
			return
		}
		rev_buf := string(buf[0 : n-1]) //delete the tail #
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

func testbuf() {
	buf := "S168#358511029674984#0007#0013#B2G:1cc,0,247c,1300"

	var arr_buf, data_buf, comand_buf, lbs_buf []string

	arr_buf = strings.Split(buf, "#")                    //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割:
	lbs_buf = strings.Split(string(comand_buf[1]), ",")  //分割,

	flag := BDYString.HexString2Int(string(lbs_buf[0]))
	fmt.Println(flag)
}

func ParseProtocol(rev_buf string, conn net.Conn) {
	var err error
	var arr_buf, data_buf, comand_buf []string

	fmt.Println("Receive from client", rev_buf)

	arr_buf = strings.Split(rev_buf, "#")                //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割;

	fmt.Println(comand_buf[0])

	switch comand_buf[0] {
	case "LOCA":
		//parse data
		imei := string(arr_buf[1])
		serial_num := string(arr_buf[2])
		switch comand_buf[1] {
		case "W":
			alert := BDYString.GetBetweenStr(rev_buf, "ALERT", ";")
			status := BDYString.GetBetweenStr(rev_buf, "STATUS", ";")
			//wifi := BDYString.GetBetweenStr(rev_buf, "WIFI", "$")
			fmt.Println(status)
			fmt.Println(alert)
			//fmt.Println(wifi)
			break
		case "G":

			break
		case "L":

			break
		}
		//printf data
		//send data
		buf := fmt.Sprintf("S168#%s#%s##ACK^LOCA,", imei, serial_num)
		buf_send := fmt.Sprintf("S168#%s#%s#%04d#ACK^LOCA,", imei, serial_num, len([]rune(buf))-27)
		fmt.Println(buf_send)
		_, err = conn.Write([]byte(buf_send))
		break
	case "B2G":
		//parse data
		var lbs_buf []string
		var lbs_int [4]int
		lbs_buf = strings.Split(string(comand_buf[1]), ",") //分割;
		for i := 0; i < 0; i++ {
			lbs_int[i] = BDYString.HexString2Int(string(lbs_buf[i]))
		}
		fmt.Println(lbs_int)
		//printf data

		//send data
		break
	}
	if err != nil {
		return
	}
}
