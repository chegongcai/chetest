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

func Int2HexString(lenght int) string {
	var num [4]int
	var buf [4]string
	num[0] = lenght / 4096
	num[1] = lenght % 4096 / 256
	num[2] = lenght % 4096 % 256 / 16
	num[3] = lenght % 16
	for i := 0; i < 4; i++ {
		if num[i] == 10 {
			buf[i] = "a"
		} else if num[i] == 11 {
			buf[i] = "b"
		} else if num[i] == 12 {
			buf[i] = "c"
		} else if num[i] == 13 {
			buf[i] = "d"
		} else if num[i] == 14 {
			buf[i] = "e"
		} else if num[i] == 15 {
			buf[i] = "f"
		} else {
			buf[i] = strconv.Itoa(num[i])
		}
	}
	str_out := buf[0] + buf[1] + buf[2] + buf[3]
	return str_out
}

func HexString2Int(str string) int {
	var num [4]int
	buf := []byte(str)
	for i := 0; i < 4; i++ {
		if string(buf[i]) == "a" {
			num[i] = 10
		} else if string(buf[i]) == "b" {
			num[i] = 11
		} else if string(buf[i]) == "c" {
			num[i] = 12
		} else if string(buf[i]) == "d" {
			num[i] = 13
		} else if string(buf[i]) == "e" {
			num[i] = 14
		} else if string(buf[i]) == "f" {
			num[i] = 15
		} else {
			num[i], _ = strconv.Atoi(string(buf[i]))
		}
	}
	flag := num[0]*4096 + num[1]*256 + num[2]*16 + num[3]
	return flag
}

func testbuf() {
	buf := fmt.Sprintf("S168#358511029674984#0011##ACK^LOCA,123")
	lenght := len([]rune(buf)) - 27
	fmt.Println(lenght)
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
		//S168#000000000000008#0f12#0071#LOCA:G;CELL:1,1cc,2,2795,1435,64;GDATA:A,12,160412154800,22.564025,113.242329,5.5,152,900;ALERT:0000;STATUS:89,98;WAY:0$
		//S168#358511029674984#0020#019d#LOCA:W;CELL:4,1cc,0,247c,1300,3b,247c,f49,12,247c,f47,16,247c,f48,17;GDATA:V,0,180719152738,0.000000,0.000000,0,0,0;ALERT:0080;STATUS:100,100;WIFI:12,bc-1c-81-2d-8f-64,-54,50-0f-f5-0b-fd-c1,-67,0a-9b-4b-98-51-91,-70,08-9b-4b-98-51-91,-71,0a-9b-4b-98-51-69,-71,08-9b-4b-98-51-69,-74,a8-ad-3d-f5-56-8c,-77,84-a9-c4-64-3e-22,-84,08-9b-4b-98-51-6d,-88,0a-9b-4b-98-51-6d,-88,84-a9-c4-4d-7e-e2,-88,08-9b-4b-98-51-59,-89$
		//parse data
		imei := string(arr_buf[1])
		serial_num := string(arr_buf[2])
		switch comand_buf[1] {
		case "W":

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
		buf_send := fmt.Sprintf("S168#%s#%s#%d#ACK^LOCA,", imei, serial_num, len([]rune(buf))-27)
		fmt.Println(buf_send)
		_, err = conn.Write([]byte(buf_send))
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
