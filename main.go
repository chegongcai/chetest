package main

import (
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

	var arr_buf, data_buf, comand_buf []string

	var flag string = "S168#000000000000008#0f12#0071#LOCA:G;CELL:1,1cc,2,2795,1435,64;GDATA:A,12,160412154800,22.564025,113.242329,5.5,152,900;ALERT:0000;STATUS:89,98;WAY:0$"

	arr_buf = strings.Split(flag, "#")                   //先分割#
	data_buf = strings.Split(string(arr_buf[4]), ";")    //分割;
	comand_buf = strings.Split(string(data_buf[0]), ":") //分割;

	fmt.Println(comand_buf[0])
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
	case "LOCA": //S168#000000000000008#0f12#0071#LOCA:G;CELL:1,1cc,2,2795,1435,64;GDATA:A,12,160412154800,22.564025,113.242329,5.5,152,900;
		//ALERT:0000;STATUS:89,98;WAY:0$
		//parse data
		imei := string(arr_buf[1])
		serial_num := string(arr_buf[2])
		//printf data
		fmt.Println("imei:", imei)
		fmt.Println("serial_num:", serial_num)
		//send data
		buf := fmt.Sprintf("S168#%s#%s#0009#ACK^LOCA,", imei, serial_num)
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT02": //BDT02,180528,22.5641N,113.2524E,000.1,061830,323.87,060009080002#
		//parse data

		//printf data

		//send data
		buf := fmt.Sprintf("BDS02#")
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT03": //BDT03,7,460,0,9520|3671|13,9520|3672|12,9520|3673|11,9520|3674|10,9520|3675|9,9520|3676|8,9520|3677|7#
		//parse data

		//printf data

		//send data
		buf := fmt.Sprintf("BDS03#")
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	case "BDT04": //BDT04,060009080002#
		//parse data

		//printf data

		//send data
		buf := fmt.Sprintf("BDS04#")
		fmt.Println(buf)
		_, err = conn.Write([]byte(buf))
		break
	}
	if err != nil {
		return
	}
}
