package main

import (
	"crypto/tls"
	"encoding/json"
	"net"

	con "github.com/jmarcantony/tchat/connection"
)

func handleMessageConnection(conn con.Connection, messageChan chan []byte) {
	id := conn.Read()
	ok, room := rooms.Exists(id)
	if !ok {
		conn.Write([]byte("1"))
		return
	}
	host, _, err := net.SplitHostPort(conn.C.RemoteAddr().String())
	if err != nil {
		log.Println(err)
	}
	ok, member := room.IsMember(host)
	if !ok {
		conn.Write([]byte("1"))
		return
	}
	conn.Write([]byte("0"))
	for {
		select {
		case msg := <-messageChan:
			var data Message
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println(err)
			}
			if data.Id == room.Id {
				conn.Write([]byte(data.Nickname + ": " + data.Message))
			}
		case <-member.Ctx.Done():
			return
		}
	}
}

func runMessageServer(serverKey, serverCert []byte) {
	cer, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	l, err := tls.Listen("tcp4", ":"+messagePort, config)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Message Server listening on port %s", messagePort)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		conn := con.Connection{C: c}
		ch := make(chan []byte, 100) // make size dynamic
		messageChanels = append(messageChanels, ch)
		go handleMessageConnection(conn, ch)
	}
}
