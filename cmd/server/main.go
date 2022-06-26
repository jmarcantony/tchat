package main

import (
	"net"

	con "github.com/jmarcantony/tchat/connection"
	"github.com/jmarcantony/tchat/logger"
	sec "github.com/jmarcantony/tchat/security"
)

const port = "9000"

var (
	log                   = logger.NewLogger("server.log")
	connections           []net.Conn
	publicKey, privateKey []byte
)

func handshake(conn net.Conn) con.Connection {
	c := con.Connection{C: conn}
	buf := make([]byte, 32)
	conn.Read(buf)
	conn.Write(publicKey)
	c.PeerKey = buf
	c.PrivateKey = privateKey

	return c
}

func handleConnection(conn con.Connection) {
	for {
		msg := conn.Read()
		_ = msg
	}
}

func main() {
	publicKey, privateKey = sec.GenerateKey()
	l, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on port %s", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		connections = append(connections, conn)
		log.Printf("Connection recieved from %s", conn.LocalAddr().String())
		go handleConnection(handshake(conn))
	}
}
