package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"time"
)

var (
	mode = flag.String("m", "server", "mode: client or server")
	port = flag.String("p", "4000", "host: ip:port")
)

func main() {
	flag.Parse()

	switch *mode {
	case "server":
		RunServer()
	case "client":
		RunClient()
	}
}

func RunServer() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+*port)
	CheckError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	log.Println("Listening on port " + *port)
	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		packet := &Packet{}
		err = proto.Unmarshal(buf[0:n], packet)
		log.Printf("Received %d sent at %s from %s", *packet.Serial, time.Unix(*packet.SentTime, 0), addr)

		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
}

func RunClient() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+*port)
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	CheckError(err)

	defer conn.Close()

	i := 1
	for {
		packet := CreatePacket(int32(i), "dummy message")
		now := time.Now().Unix()
		packet.SentTime = &now
		data, err := proto.Marshal(packet)
		if err != nil {
			log.Fatal("marshalling error: ", err)
		}
		buf := []byte(data)
		_, err = conn.Write(buf)
		if err != nil {
			log.Println(err)
		}

		i++
		time.Sleep(time.Second * 1)
	}

}

func CreatePacket(serial int32, msg string) *Packet {
	packet := Packet{
		Serial:  &serial,
		Message: &msg,
	}
	return &packet
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
