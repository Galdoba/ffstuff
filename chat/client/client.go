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
}

func New(conn net.Conn) Client {
	return Client{
		conn:     conn,
		nick:     "Anonymous",
		room:     &protocol.Room{},
		commands: make(chan<- protocol.Command),
	}
}

func (c *Client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: '%s'", cmd))
		}
	}
}

func (c *Client) err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func (c *Client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *Client) Conn() net.Conn {
	return c.conn
}
