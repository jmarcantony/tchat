package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"

	con "github.com/jmarcantony/tchat/connection"
	"github.com/jmarcantony/tchat/logger"
)

const port = "9000"

var (
	log         = logger.NewLogger("server.log")
	connections []net.Conn
)

func handleConnection(conn con.Connection) {
	for {
		msg := conn.Read()
		fmt.Println(msg)
	}
}

func loadCertKey() ([]byte, []byte) {
	key, err := ioutil.ReadFile("localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	cert, err := ioutil.ReadFile("localhost.pem")
	if err != nil {
		log.Fatal(err)
	}
	return key, cert
}

func main() {
	serverKey, serverCert := loadCertKey()
	cer, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	l, err := tls.Listen("tcp4", ":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Server listening on port %s", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		connections = append(connections, conn)
		log.Printf("Connection recieved from %s", conn.LocalAddr().String())
		go handleConnection(con.Connection{C: conn})
	}
}
