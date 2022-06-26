package connection

import (
	"bytes"
	"log"
	"net"
	"strconv"
)

type Connection struct {
	C                   net.Conn
	PeerKey, PrivateKey []byte
}

func (c Connection) GetMsgSize() int {
	buf := make([]byte, 1000)
	c.C.Read(buf)
	n, _ := strconv.Atoi(string(buf))
	return n
}

func (c Connection) Read() string {
	l := c.GetMsgSize()
	buf := make([]byte, l)
	n, err := c.C.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	encrypted := buf[:n]
	_ = encrypted
	// TODO: decrypt and return string
	return ""
}
func (c Connection) Write(msg []byte) {
	m := bytes.TrimSpace(msg)
	c.C.Write([]byte(strconv.Itoa(len(m))))
	c.C.Write(m)
}
