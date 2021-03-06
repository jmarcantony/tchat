package room

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	con "github.com/jmarcantony/tchat/connection"
)

type Room struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Len      int       `json:"len"`
	Private  bool      `json:"private"`
	Members  []*Member `json:"-"`
}

func (r Room) IsMember(remoteAddr string) (bool, *Member) {
	for _, member := range r.Members {
		host, _, err := net.SplitHostPort(member.Conn.C.RemoteAddr().String())
		if err != nil {
			log.Println(err)
		}
		if host == remoteAddr {
			return true, member
		}
	}
	return false, nil
}

type RoomJson struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Len  int    `json:"len"`
}

type RoomStatus struct {
	Status   string `json:"s"`
	Password string `json:"p"`
	Name     string `json:"n"`
}

type Member struct {
	Nickname string
	Conn     con.Connection
	Admin    bool
	Ctx      context.Context
	Cancel   context.CancelFunc
}

type Rooms []*Room

func (r Rooms) Public() Rooms {
	var publicRooms Rooms
	for _, room := range r {
		if !room.Private {
			publicRooms = append(publicRooms, room)
		}
	}
	return publicRooms
}

func (r Rooms) Exists(id string) (bool, *Room) {
	for i, j := 0, len(r)-1; i <= j; i, j = i+1, j-1 {
		if r[i].Id == id {
			return true, r[i]
		} else if r[j].Id == id {
			return true, r[j]
		}
	}
	return false, nil
}

func NewRoom(name string, private bool, idLen int) (*Room, error) {
	id := make([]byte, idLen)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	return &Room{Name: name, Id: fmt.Sprintf("%X", id), Private: private}, nil
}
