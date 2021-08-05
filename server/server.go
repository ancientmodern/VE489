package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	. "ve489/util"
)

const (
	IP       = "0.0.0.0"
	Filepath = "/root/VE489/received_text.txt"
)

var port = flag.Int("p", 8002, "Server Port")

func main() {
	file, err := os.OpenFile(Filepath, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("OpenFile error", err)
	}
	defer file.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", IP, *port))
	if err != nil {
		fmt.Println("ResolveUDPAddr err:", err)
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("ListenUDP err:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening on port", *port)

	rxSeqNum, count := false, 0
	for {
		str := ""
		buf := make([]byte, 1024)
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("ReadFromUDP err:", err)
			return
		}

		txSeqNum := Byte2Bool(buf[n-1])

		if txSeqNum == rxSeqNum {
			fmt.Printf("Want %d, received %d, send back ACK %d\n", Bool2Int(rxSeqNum), Bool2Int(txSeqNum), Bool2Int(!rxSeqNum))
			rxSeqNum = !rxSeqNum
			count++
			str += string(buf[:n-1])
			fmt.Println("Totally received", count)
		} else {
			fmt.Printf("Want %d, received %d, send back ACK %d. Drop this message\n", Bool2Int(rxSeqNum), Bool2Int(txSeqNum), Bool2Int(rxSeqNum))
		}

		_, err = conn.WriteToUDP([]byte(fmt.Sprintf("ACK %d", Bool2Int(rxSeqNum))), cliAddr)
		if err != nil {
			fmt.Println("WriteToUDP err:", err)
			return
		}

		buf = nil

		write := bufio.NewWriter(file)
		if _, err = write.WriteString(str); err != nil {
			return
		}
		write.Flush()
	}

	//err = ioutil.WriteFile("/root/VE489/received_text.txt", []byte(s), 0777)
	//if err != nil {
	//	fmt.Println("WriteFile error:", err)
	//	return
	//}
}
