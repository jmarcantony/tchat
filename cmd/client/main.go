package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	con "github.com/jmarcantony/tchat/connection"
	sec "github.com/jmarcantony/tchat/security"
)

const (
	serverAddr = "localhost"
	port       = "9000"
)

var currentRoom string

func handshake(c net.Conn) con.Connection {
	conn := con.Connection{C: c}
	pubKey, privKey := sec.GenerateKey()
	conn.PrivateKey = privKey
	conn.C.Write(pubKey)

	return conn
}

func main() {
	c, err := net.Dial("tcp4", serverAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	conn := handshake(c)
	s := bufio.NewScanner(os.Stdin)
	fmt.Println(`
████████╗ ██████╗██╗  ██╗ █████╗ ████████╗
╚══██╔══╝██╔════╝██║  ██║██╔══██╗╚══██╔══╝
   ██║   ██║     ███████║███████║   ██║   
   ██║   ██║     ██╔══██║██╔══██║   ██║   
   ██║   ╚██████╗██║  ██║██║  ██║   ██║   
   ╚═╝    ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝  
                                        (@jmarcantony)
	`)
	for {
		fmt.Printf("tchat%s>> ", currentRoom)
		s.Scan()
		conn.Write(s.Bytes())
	}
}
