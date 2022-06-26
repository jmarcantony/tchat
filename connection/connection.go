package connection

import (
	"bytes"
	"log"
	"net"
	"strconv"
)

type Connection struct {
	C net.Conn
}

func (c Connection) GetMsgSize() int {
	buf := make([]byte, 1000)
	n, err := c.C.Read(buf)
	if err != nil {
		log.Println(err)
	}
	l, _ := strconv.Atoi(string(buf[:n]))
	return l
}

func (c Connection) Read() string {
	l := c.GetMsgSize()
	buf := make([]byte, l)
	c.Proceed('0')
	n, err := c.C.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	msg := buf[:n]
	return string(msg)
}

func (c Connection) Write(msg []byte) {
	m := bytes.TrimSpace(msg)
	c.C.Write([]byte(strconv.Itoa(len(m))))
	<-c.Await()
	c.C.Write(m)
}

func (c Connection) Proceed(code byte) {
	c.C.Write([]byte{code})
}

func (c Connection) Await() chan []byte {
	ch := make(chan []byte, 1)
	buf := make([]byte, 1)
	go func() {
		c.C.Read(buf)
		ch <- buf
	}()
	return ch
}
