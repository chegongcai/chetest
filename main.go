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
	buf := "S168#358511029674984#0020#019d#LOCA:W;CELL:4,1cc,0,247c,1300,3b,247c,f49,12,247c,f47,16,247c,f48,17;GDATA:V,0,180719152738,0.000000,0.000000,0,0,0;ALERT:0080;STATUS:100,100;WIFI:12,bc-1c-81-2d-8f-64,-54,50-0f-f5-0b-fd-c1,-67,0a-9b-4b-98-51-91,-70,08-9b-4b-98-51-91,-71,0a-9b-4b-98-51-69,-71,08-9b-4b-98-51-69,-74,a8-ad-3d-f5-56-8c,-77,84-a9-c4-64-3e-22,-84,08-9b-4b-98-51-6d,-88,0a-9b-4b-98-51-6d,-88,84-a9-c4-4d-7e-e2,-88,08-9b-4b-98-51-59,-89"
	//var arr_buf, data_buf, comand_buf, w_buf []string

	//arr_buf = strings.Split(buf, "#")                    //先分割#
	//data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	//comand_buf = strings.Split(string(data_buf[0]), ":") //分割;

	alert := BDYString.GetBetweenStr(buf, "ALERT", ";")
	wifi := BDYString.GetBetweenStr(buf, "WIFI", "$")
	fmt.Println(alert)
	fmt.Println(wifi)
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
			wifi := BDYString.GetBetweenStr(rev_buf, "WIFI", "$")
			fmt.Println(status)
			fmt.Println(alert)
			fmt.Println(wifi)
			break
		case "G":

			break
		case "L":

			break
		}
		//printf data
		fmt.Println("imei:", imei)
		fmt.Println("serial_num:", serial_num)
		//send data
		buf := fmt.Sprintf("S168#%s#%s##ACK^LOCA,", imei, serial_num)
		buf_send := fmt.Sprintf("S168#%s#%s#%04d#ACK^LOCA,", imei, serial_num, len([]rune(buf))-27)
		fmt.Println(buf_send)
		_, err = conn.Write([]byte(buf_send))
		break
	case "BDT02":
		//parse data

		//printf data

		//send data
		break
	case "BDT03":
		//parse data

		//printf data

		//send data
		break
	case "BDT04":
		//parse data

		//printf data

		//send data

		break
	}
	if err != nil {
		return
	}
}
