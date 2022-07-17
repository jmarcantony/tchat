package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	con "github.com/jmarcantony/tchat/connection"
	"github.com/jmarcantony/tchat/logger"
	"github.com/jmarcantony/tchat/room"
)

// TODO: Use flags and yaml file to load port data etc
const (
	port        = "9000"
	messagePort = "9090"
)

var (
	log            = logger.NewLogger("server.log")
	connections    []net.Conn
	rooms          room.Rooms
	messageChanels []chan []byte // make size dynamic in config
)

type Message struct {
	Id       string `json="i"`
	Nickname string `json="n"`
	Message  string `json="m"`
}

func handleRoom(conn con.Connection, r *room.Room) {
	// TODO: Handle unique nicknames
	// TODO: Handle File upload and downloads
	// TODO: Handle Personal Messages
	// TODO: Handle Banning and Kicking
	r.Len++
	ctx, cancel := context.WithCancel(context.Background())
	// TODO: cancel context when member is no longer in room
	member := &room.Member{Conn: conn, Nickname: "Anonymous", Ctx: ctx, Cancel: cancel} // TODO: Read Nickname
	r.Members = append(r.Members, member)
	// TODO: Log joining of new member
	fmt.Println("Someone joined", r.Name)
	for {
		text := conn.Read()
		switch cmd := strings.Split(text, ":"); cmd[0] {
		case "n": // n short for nickname
			fmt.Println("someone changed nickname to", cmd[1])
			member.Nickname = cmd[1]
		default:
			m := Message{Id: r.Id, Nickname: member.Nickname, Message: text[2:]}
			data, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
			}
			for _, messageChan := range messageChanels {
				messageChan <- data
			}
		}
	}
}

func handleConnection(conn con.Connection) {
	for {
		msg := conn.Read()
		if len(msg) == 0 {
			continue
		}
		switch cmd := strings.Split(msg, " "); cmd[0] {
		case "r":
			data, _ := json.Marshal(rooms.Public())
			conn.Write(data)
		case "j":
			id := cmd[1]
			ok, room := rooms.Exists(id)
			if !ok {
				errJson, _ := json.Marshal(map[string]string{"s": "1", "p": "", "n": ""})
				conn.Write(errJson)
			} else {
				if room.Private {
					data, _ := json.Marshal(map[string]string{"s": "0", "p": room.Password, "n": room.Name})
					conn.Write(data)
				} else {
					data, _ := json.Marshal(map[string]string{"s": "0", "p": "", "n": room.Name})
					conn.Write(data)
				}
				msg := conn.Read()
				if msg == "0" {
					handleRoom(conn, room)
				}
			}
		}
		// TODO: Handle State (in a room or not) and query commands from client
		// TODO: Filter server command symbols
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
	go runMessageServer(serverKey, serverCert)
	generalRoom, err := room.NewRoom("general", false, 5)
	if err != nil {
		log.Fatal(err)
	}
	rooms = append(rooms, generalRoom)
	randomRoom, err := room.NewRoom("random", false, 5)
	if err != nil {
		log.Fatal(err)
	}
	rooms = append(rooms, randomRoom)
	privRoom, err := room.NewRoom("private room", true, 5)
	if err != nil {
		log.Fatal(err)
	}
	rooms = append(rooms, privRoom)
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
