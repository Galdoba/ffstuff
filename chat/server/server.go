package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Galdoba/ffstuff/chat/client"
	"github.com/Galdoba/ffstuff/chat/protocol"
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

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.ID {
		case protocol.CMD_NICK:
			s.cmd_nick(cmd.Sender, cmd.Args)
		case protocol.CMD_JOIN:
			s.cmd_join(cmd.Sender, cmd.Args)
		case protocol.CMD_ROOMS:
			s.cmd_rooms(cmd.Sender, cmd.Args)
		case protocol.CMD_MSG:
			s.cmd_msg(cmd.Sender, cmd.Args)
		case protocol.CMD_QUIT:
			s.cmd_quit(cmd.Sender, cmd.Args)

		}
	}
}

func (s *server) NewClient(conn net.Conn) {
	log.Printf("new client has connected: %v", conn.RemoteAddr().String())
	// c := &client{
	// 	conn:     conn,
	// 	nick:     "Anonymos",
	// 	room:     &room{},
	// 	commands: s.commands,
	// }
	c := client.New(conn)

	c.ReadInput()
}

func (s *server) cmd_nick(c *client, args []string) error {
	c.nick = args[1]
	c.msg(fmt.Sprintf("You are now: %s", c.nick))
	return nil
}

func (s *server) cmd_join(c *client, args []string) error {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.leaveCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("welcome to %s", r.name))
	return nil
}

func (s *server) cmd_rooms(c *client, args []string) error {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are: [%s]", strings.Join(rooms, ", ")))
	return nil
}

func (s *server) cmd_msg(c *client, args []string) error {
	if c.room == nil {
		c.err(fmt.Errorf("you must join the room first"))
		return nil
	}
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
	return nil
}

func (s *server) cmd_quit(c *client, args []string) error {
	log.Printf("client has disconnected: %s (%s)", c.nick, c.conn.RemoteAddr().String())

	s.leaveCurrentRoom(c)

	c.msg("have a good day!")

	c.conn.Close()

	return nil
}

func (s *server) leaveCurrentRoom(c *client) error {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
	return nil
}
