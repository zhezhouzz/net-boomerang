package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var CLIENT_FILE_PATH string = "../data/clinet_data.txt"
var BEGIN_PATTERN string = "start-->"
var END_PATTERN string = "<--end"
var SERVER_PORT string = ":8888"

func clientRead(conn net.Conn) int {
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

func clientWrite(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Fatalf("send [%s] content faild: %s\n", string(data), err)
	}
	log.Printf("send [%s] content success\n", string(data))
}

func clientRecv(conn net.Conn, data []byte) {
	_, err := conn.Read(data)
	if err != nil {
		log.Fatalf("read content from conn failed\n")
	}
	log.Printf("recv [%s] content success\n", string(data))
}

func clientConn(conn net.Conn) {
	defer conn.Close()

	clientWrite(conn, []byte(BEGIN_PATTERN))
	off := clientRead(conn)

	fp, err := os.OpenFile(CLIENT_FILE_PATH, os.O_RDONLY, 0755)
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
				clientWrite(conn, []byte(END_PATTERN))
				log.Println("send all content, now wait for sendback")
				for {
					data := make([]byte, 10)
					clientRecv(conn, data)
					if err != nil {
						continue
					} else {
						log.Printf("recv finished, conn now close\n")
						break
					}
				}
				break
			}
			log.Fatalf("read file err: %s\n", err)
		}
		clientWrite(conn, data[:n])
	}
}

func main() {
	conn, err := net.DialTimeout("tcp", SERVER_PORT, time.Second*10)
	if err != nil {
		log.Fatalf("client dial faild: %s\n", err)
	}
	clientConn(conn)
}
