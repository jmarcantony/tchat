package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	con "github.com/jmarcantony/tchat/connection"
)

const (
	serverAddr = "localhost"
	port       = "9000"
)

var currentRoom string

func loadCert() []byte {
	c, err := ioutil.ReadFile("../server/localhost.pem")
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func main() {
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(loadCert()); !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots, ServerName: "localhost"}
	c, err := tls.Dial("tcp4", serverAddr+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	conn := con.Connection{C: c}
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
