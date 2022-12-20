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
		cm, _ := protocol.Assemble(cmd)
		cm.Sender = c
		go func() { c.commands <- cm }()
		/*switch cmd {
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
		}*/
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

func (c *Client) LeaveRoom() {
	c.room = nil
}
