package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var SENDER_FILE_PATH string = "../data/sender_data.txt"
var BEGIN_PATTERN string = "start-->"
var END_PATTERN string = "<--end"
var SENDER_PORT string = ":8888"

func senderRead(conn net.Conn) int {
	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("receive server info faild: %s\n", err)
	}
	// string conver int
	off, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalf("string conver int faild: %s\n", err)
	}
	return off
}

func senderWrite(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Fatalf("send [%s] content faild: %s\n", string(data), err)
	}
	log.Printf("send [%s] content success\n", string(data))
}

func senderRecv(conn net.Conn, data []byte) {
	_, err := conn.Read(data)
	if err != nil {
		log.Fatalf("read content from conn failed\n")
	}
	log.Printf("recv [%s] content success\n", string(data))
}

func requestHandle(conn net.Conn) {
	defer conn.Close()

	senderWrite(conn, []byte(BEGIN_PATTERN))
	off := senderRead(conn)

	fp, err := os.OpenFile(SENDER_FILE_PATH, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalf("open file faild: %s\n", err)
	}
	defer fp.Close()

	_, err = fp.Seek(int64(off), 0)
	if err != nil {
		log.Fatalf("set file seek faild: %s\n", err)
	}
	log.Printf("read file at seek: %d\n", off)

	for {
		data := make([]byte, 10)
		n, err := fp.Read(data)
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Second * 1)
				senderWrite(conn, []byte(END_PATTERN))
				log.Println("send all content, now wait for sendback")
				break
			}
			log.Fatalf("read file err: %s\n", err)
		}
		senderWrite(conn, data[:n])
	}
}

func main() {
	l, err := net.Listen("tcp", SENDER_PORT)
	if err != nil {
		log.Fatalf("error listen: %s\n", err)
	}
	defer l.Close()

	for {
		log.Println("waiting accept.")
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("accept faild: %s\n", err)
		}
		requestHandle(conn)
	}
}
