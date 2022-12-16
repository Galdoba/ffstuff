package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ffstuff/chat/client"
	"github.com/ffstuff/chat/protocol"
)

type server struct {
	rooms    map[string]*protocol.Room
	commands chan protocol.Command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*protocol.Room),
		commands: make(chan protocol.Command),
	}
}

func (s *server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case protocol.CMD_NICK:
			s.cmd_nick(cmd.Client, cmd.Args)
		case protocol.CMD_JOIN:
			s.cmd_join(cmd.Client, cmd.Args)
		case protocol.CMD_ROOMS:
			s.cmd_rooms(cmd.Client, cmd.Args)
		case protocol.CMD_MSG:
			s.cmd_msg(cmd.Client, cmd.Args)
		case protocol.CMD_QUIT:
			s.cmd_quit(cmd.Client, cmd.Args)

		}
	}
}

func (s *server) NewClient(conn net.Conn) {
	log.Printf("new client has connected: %v", conn.RemoteAddr().String())
	c := &client.Client{
		Conn:     conn,
		Nick:     "Anonymos",
		Room:     &protocol.Room{},
		Commands: s.commands,
	}

	c.ReadInput()
}

func (s *server) cmd_nick(c *client.Client, args []string) error {
	c.Nick = args[1]
	c.Msg(fmt.Sprintf("You are now: %s", c.Nick))
	return nil
}

func (s *server) cmd_join(c *client.Client, args []string) error {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &protocol.Room{
			Name:    roomName,
			Members: make(map[net.Addr]*client.Client),
		}
		s.rooms[roomName] = r
	}

	r.Members[c.Conn.RemoteAddr()] = c

	s.leaveCurrentRoom(c)

	c.Room = r

	r.Broadcast(c, fmt.Sprintf("%s has joined the room", c.Nick))
	c.Msg(fmt.Sprintf("welcome to %s", r.Name))
	return nil
}

func (s *server) cmd_rooms(c *client.Client, args []string) error {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.Msg(fmt.Sprintf("available rooms are: [%s]", strings.Join(rooms, ", ")))
	return nil
}

func (s *server) cmd_msg(c *client.Client, args []string) error {
	if c.Room == nil {
		c.Err(fmt.Errorf("you must join the room first"))
		return nil
	}
	c.Room.Broadcast(c, c.Nick+": "+strings.Join(args[1:], " "))
	return nil
}

func (s *server) cmd_quit(c *client.Client, args []string) error {
	log.Printf("client has disconnected: %s (%s)", c.Nick, c.Conn.RemoteAddr().String())

	s.leaveCurrentRoom(c)

	c.Msg("have a good day!")

	c.Conn.Close()

	return nil
}

func (s *server) leaveCurrentRoom(c *client.Client) error {
	if c.Room != nil {
		delete(c.Room.Members, c.Conn.RemoteAddr())
		c.Room.Broadcast(c, fmt.Sprintf("%s has left the room", c.Nick))
	}
	return nil
}
