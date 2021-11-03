package main

import (
	mynet "Samurai/net"
	"io"
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
		pack := mynet.NewPacket()
		msg, _ := pack.Pack(mynet.NewMsgPacket(0, []byte("Hi!")))
		_, err := conn.Write(msg)
		if err != nil {
			log.Println("Write error: ", err)
			return
		}
		headerData := make([]byte, pack.GetHeaderLen())
		_, err = io.ReadFull(conn, headerData)
		if err != nil {
			log.Println("Read headerData error : ", err)
			break
		}
		msgHeader, err := pack.Unpack(headerData)
		if err != nil {
			log.Println("Server unpack error : ", err)
			return
		}
		if msgHeader.GetLen() > 0 {
			msg := msgHeader.(*mynet.Message)
			msg.Data = make([]byte, msg.GetLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				log.Println("Server unpack data error : ", err)
				return
			}
			log.Printf("==>ID=%d,Len=%d,Data=%s", msg.ID, msg.Len, string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
