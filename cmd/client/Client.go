package main

import (
	gnet "gnet/net"
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
		pack := gnet.NewPacket()
		// write data
		msg, _ := pack.Pack(gnet.NewMsgPacket(0, []byte("Hi!")))
		_, err := conn.Write(msg)
		if err != nil {
			log.Println("Write error: ", err)
			return
		}
		// read data
		headerData := make([]byte, pack.GetHeaderLen())
		// read header data
		_, err = io.ReadFull(conn, headerData)
		if err != nil {
			log.Println("Read headerData error : ", err)
			break
		}
		// unpack the header
		msgHeader, err := pack.Unpack(headerData)
		if err != nil {
			log.Println("Server unpack error : ", err)
			return
		}
		// read the rest of data
		if msgHeader.GetLen() > 0 {
			msg := msgHeader.(*gnet.Message)
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
