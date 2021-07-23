package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	. "ve489/util"
)

var num = flag.Int("n", 20, "Input how many times")
var ip = flag.String("i", "10.3.39.2", "Server IP")
var port = flag.Int("p", 8002, "Server Port")

func main() {
	flag.Parse()

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Dial complete")

	txSeqNum, count := false, 0
	for i := 0; i < *num; i++ {
		_, err = conn.Write([]byte(fmt.Sprintf("Message%d", Bool2Int(txSeqNum))))
		if err != nil {
			fmt.Println("conn.Write err:", err)
		}
		fmt.Printf("Sent Message %d (Seq: %d)\n", count, Bool2Int(txSeqNum))

		err = conn.SetReadDeadline(time.Now().Add(1500 * time.Millisecond))
		if err != nil {
			fmt.Println("conn.SetReadDeadline err:", err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Printf("Waiting for ACK timeout, will resend Message %d (Seq: %d)\n", count, Bool2Int(txSeqNum))
			} else {
				return
			}
		} else {
			fmt.Println("Received from Server:", string(buf[:n]))
			count++
			txSeqNum = !txSeqNum
		}

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Client exits")
}
