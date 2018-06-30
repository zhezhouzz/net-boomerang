package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	// "time"
)

var BEGIN_PATTERN string = "start-->"
var END_PATTERN string = "<--end"
var SERVER_PORT string = ":8888"
var SERVER_FILE_PATH string = "../data/server_data.txt"

func writeFile(content []byte) {
	if len(content) != 0 {
		fp, err := os.OpenFile(SERVER_FILE_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		defer fp.Close()
		if err != nil {
			log.Fatalf("open file faild: %s\n", err)
		}
		_, err = fp.Write(content)
		if err != nil {
			log.Fatalf("append content to file faild: %s\n", err)
		}
		log.Printf("append content: [%s] success\n", string(content))
	}
}

func getFileStat() int64 {
	fileinfo, err := os.Stat(SERVER_FILE_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("file size: %d\n", 0)
			return int64(0)
		}
		log.Fatalf("get file stat faild: %s\n", err)
	}
	log.Printf("file size: %d\n", fileinfo.Size())
	return fileinfo.Size()
}

func serverConn(conn net.Conn) {
	defer conn.Close()
	for {
		var buf = make([]byte, 10)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("server io EOF\n")
			}
			log.Fatalf("server read faild: %s\n", err)
		}
		log.Printf("recevice %d bytes, content is [%s]\n", n, string(buf[:n]))
		switch string(buf[:n]) {
		case BEGIN_PATTERN:
			off := getFileStat()
			// int conver string
			stringoff := strconv.FormatInt(off, 10)
			_, err = conn.Write([]byte(stringoff))
			if err != nil {
				log.Fatalf("server write faild: %s\n", err)
			}
			continue
		case END_PATTERN:
			stringback := "hello world"
			log.Printf("receive over\n")
			_, err = conn.Write([]byte(stringback))
			if err != nil {
				log.Fatalf("server write faild: %s\n", err)
			}
			log.Fatalf("server send back finished\n")
		}
		writeFile(buf[:n])
	}
}

func main() {
	l, err := net.Listen("tcp", SERVER_PORT)
	if err != nil {
		log.Fatalf("error listen: %s\n", err)
	}
	defer l.Close()

	log.Println("waiting accept.")
	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("accept faild: %s\n", err)
	}
	serverConn(conn)
}
