package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/Galdoba/ffstuff/chat/protocol"
)

type Client struct {
	conn     net.Conn
	nick     string
	room     *protocol.Room
	commands chan<- protocol.Command
	//commandStr chan<- string
}

func New(conn net.Conn, comChan chan<- protocol.Command) *Client {
	return &Client{
		conn:     conn,
		nick:     "Anonymous",
		room:     &protocol.Room{},
		commands: comChan,
	}
}

func (c *Client) ReadInput() {
	c.Msg("you are in")
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0]) + " " + c.conn.RemoteAddr().Network() + " " + c.conn.RemoteAddr().String() + " " + strings.Join(args[1:], " ")
		fmt.Println("SEND:", cmd)
		c.Msg("SEND: " + cmd)
		cm, err := protocol.Assemble(cmd)
		if err != nil {
			c.Err(err)
			continue
		}

		cm.Sender = c
		c.commands <- cm

		fmt.Println("End CYCLE client")
	}
}

func (c *Client) Err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func (c *Client) Msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *Client) Conn() net.Conn {
	return c.conn
}
func (c *Client) Nick() string {
	return c.nick
}

func (c *Client) SetNick(n string) {
	c.nick = n
}

func (c *Client) Room() string {
	if c.room == nil {
		fmt.Println("------NO ROOM")
		return "[NO ROOM]"
	}
	return c.room.Name()
}

func (c *Client) CurrentRoom() *protocol.Room {
	return c.room
}

func (c *Client) LeaveRoom() {
	c.room = nil
}

func (c *Client) DebugMessage() {
	s := fmt.Sprintf("%v", c)
	c.Msg(s)
}
