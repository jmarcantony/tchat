package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/inancgumus/screen"
	"github.com/jedib0t/go-pretty/v6/table"
	con "github.com/jmarcantony/tchat/connection"
	"github.com/jmarcantony/tchat/room"
)

const (
	serverAddr  = "localhost"
	port        = "9000"
	messagePort = "9090"
)

var (
	currentRoom string
	help        = ``
	t           = table.NewWriter()
)

func loadCert() []byte {
	c, err := ioutil.ReadFile("localhost.pem")
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func recieveMessages(conn con.Connection, id string) {
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(loadCert()); !ok {
		log.Fatal("failed to parse root certificate")
	}
	config := &tls.Config{RootCAs: roots, ServerName: "localhost"}
	cReg, err := tls.Dial("tcp4", serverAddr+":"+messagePort, config)
	if err != nil {
		log.Fatal(err)
	}
	c := con.Connection{C: cReg}
	c.Write([]byte(id))
	status := c.Read()
	if status != "0" {
		return
	}
	for {
		fmt.Println("\n" + c.Read())
		fmt.Printf("tchat@%s>> ", currentRoom)
	}
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
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Members"})
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
		switch cmd := strings.Split(s.Text(), " "); cmd[0] {
		case "help":
			fmt.Println(help)
		case "cls", "clear":
			screen.Clear()
		case "rooms":
			conn.Write([]byte("r"))
			roomsString := conn.Read()
			var rooms []room.RoomJson
			if err := json.Unmarshal([]byte(roomsString), &rooms); err != nil {
				log.Println(err)
			}
			var rows []table.Row
			for _, room := range rooms {
				rows = append(rows, table.Row{room.Id, room.Name, room.Len})
			}
			t.AppendRows(rows)
			t.Render()
			t.ResetRows()
		case "join":
			if len(cmd) < 2 {
				fmt.Println("No id given")
				break
			}
			conn.Write([]byte("j " + cmd[1]))
			dataString := conn.Read()
			var status room.RoomStatus
			if err := json.Unmarshal([]byte(dataString), &status); err != nil {
				log.Println(err)
			}
			if status.Status == "1" {
				fmt.Println("Id does not exist")
			} else {
				if status.Password != "" {
					// TODO: Ask for password
					fmt.Printf("Enter password: ") // TODO: Compare hashed passwords
					s.Scan()
					if s.Text() != status.Password {
						conn.Write([]byte("1"))
						break
					}
				}
				conn.Write([]byte("0"))
				currentRoom = status.Name
				go recieveMessages(conn, cmd[1])
				fmt.Printf("tchat@%s>> ", currentRoom)
				for {
					// TODO: Handle Room Functions
					// TODO: Set nickname
					s.Scan()
					text := s.Text()
					switch cmd := strings.Split(text, " "); cmd[0] {
					case "/nick":
						fmt.Printf("Enter Nickname: ")
						s.Scan()
						nickname := s.Text() // TODO: Filter nickname for symbols and space
						if nickname != "" {
							conn.Write([]byte("n:" + nickname))
						}
						fmt.Printf("tchat@%s>> ", currentRoom)
					default:
						// TODO: Broadcast text to all memebers
						if strings.TrimSpace(text) != "" {
							conn.Write([]byte("m:" + text))
						} else {
							fmt.Printf("tchat@%s>> ", currentRoom)
						}
					}
				}
			}
		default:
			fmt.Println("Invalid command, get some 'help'")
		}
	}
}
