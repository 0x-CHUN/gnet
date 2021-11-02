package main

import (
	"log"
	"net"
	"time"
)

func main() {
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalln(err)
		return
	}
	for {
		_, err := conn.Write([]byte("hi"))
		if err != nil {
			log.Fatalln("Write err:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Read err:", err)
			return
		}
		log.Printf("Server call back : %s, cnt = %d", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
