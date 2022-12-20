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

func (s *server) Run() {
	fmt.Println("SERVER RUN")
	fmt.Println(s.commands)
	fmt.Println("-----------")
	for cmd := range s.commands {
		fmt.Println("new command:", cmd)
		switch cmd.ID {
		case protocol.CMD_NICK:
			fmt.Println(0, s)
			s.cmd_nick(cmd.Sender, cmd.Args)
			fmt.Println(1, s)
		case protocol.CMD_JOIN:
			fmt.Println(2, s)
			s.cmd_join(cmd.Sender, cmd.Args)
		case protocol.CMD_ROOMS:
			s.cmd_rooms(cmd.Sender, cmd.Args)
		case protocol.CMD_MSG:
			s.cmd_msg(cmd.Sender, cmd.Args)
		case protocol.CMD_QUIT:
			s.cmd_quit(cmd.Sender, cmd.Args)
		default:
			fmt.Println("UNK", cmd)
		}
	}
	fmt.Println("SERVER END")
}

func (s *server) NewClient(conn net.Conn) {
	log.Printf("new client has connected: %v", conn.RemoteAddr().String())
	// c := &client{
	// 	conn:     conn,
	// 	nick:     "Anonymos",
	// 	room:     &room{},
	// 	commands: s.commands,
	// }
	c := client.New(conn, s.commands)

	c.ReadInput()
	fmt.Println("END NewClient")
}

func (s *server) cmd_nick(c protocol.Receiver, args []string) error {
	//c.nick = args[1]
	fmt.Println("cmd_nick")
	switch c.(type) {
	case *client.Client:
		v := c.(*client.Client)
		fmt.Println(args)
		v.SetNick(args[0])
		v.Msg(fmt.Sprintf("You are now: %s", args[0]))
	}

	return nil
}

func (s *server) cmd_join(v protocol.Receiver, args []string) error {
	switch v.(type) {
	case *client.Client:
		c := v.(*client.Client)
		roomName := args[0]
		r, ok := s.rooms[roomName]
		if !ok {
			r = protocol.NewRoom(roomName)
			s.rooms[roomName] = r
			log.Printf("room [%v] was created", roomName)
		}
		if err := r.Join(c); err != nil {
			return fmt.Errorf("cann't execute command [/join]: %v", err.Error())
		}

		r.Broadcast(c, fmt.Sprintf("%s has joined the room", c.Nick()))
		c.Msg(fmt.Sprintf("welcome to %s", r.Name()))
	}
	return nil
}

func (s *server) cmd_rooms(v protocol.Receiver, args []string) error {
	switch v.(type) {
	case *client.Client:
		c := v.(*client.Client)
		var rooms []string
		for name := range s.rooms {
			rooms = append(rooms, name)
		}
		c.Msg(fmt.Sprintf("available rooms are: [%s]", strings.Join(rooms, ", ")))
	}
	return nil
}

func (s *server) cmd_msg(v protocol.Receiver, args []string) error {
	switch v.(type) {
	case *client.Client:
		c := v.(*client.Client)
		if c.Room() == "[NO ROOM]" {
			err := fmt.Errorf("you must join the room first")
			c.Err(err)
			return err
		}
		s.rooms[c.Room()].Broadcast(c, c.Nick()+": "+strings.Join(args, " "))
		//c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
	}
	return nil
}

func (s *server) cmd_quit(v protocol.Receiver, args []string) error {
	switch v.(type) {
	case *client.Client:
		c := v.(*client.Client)
		log.Printf("client has disconnected: %s (%s)", c.Nick(), c.Conn().RemoteAddr().String())

		s.leaveCurrentRoom(c)

		c.Msg("have a good day!")

		c.Conn().Close()
	}
	return nil
}

func (s *server) leaveCurrentRoom(c *client.Client) error {
	if c.Room() == "[NO ROOM]" {
		return fmt.Errorf("cann't leave room: client have no room to leave")
	}
	roomName := c.Room()
	s.rooms[roomName].Remove(c)
	switch len(s.rooms[roomName].Members()) {
	case 0:
		delete(s.rooms, roomName)
		log.Printf("room [%v] was abandoned", roomName)
	default:
		s.rooms[roomName].Broadcast(c, fmt.Sprintf("%s has left the room", c.Nick()))
	}
	return nil
}
