package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/ffstuff/chat/protocol"
)

type client struct {
	Conn     net.Conn
	Nick     string
	Room     *protocol.Room
	Commands chan<- protocol.Command
}

func (c *client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.Commands <- protocol.Command{
				ID:     protocol.CMD_NICK,
				Client: c,
				Args:   args,
			}
		case "/join":
			c.Commands <- protocol.Command{
				ID:     protocol.CMD_JOIN,
				Client: c,
				Args:   args,
			}
		case "/rooms":
			c.Commands <- protocol.Command{
				ID:     protocol.CMD_ROOMS,
				Client: c,
				Args:   args,
			}
		case "/msg":
			c.Commands <- protocol.Command{
				ID:     protocol.CMD_MSG,
				Client: c,
				Args:   args,
			}
		case "/quit":
			c.Commands <- protocol.Command{
				ID:     protocol.CMD_QUIT,
				Client: c,
				Args:   args,
			}
		default:
			c.Err(fmt.Errorf("unknown command: '%s'", cmd))
		}
	}
}

func (c *client) Err(err error) {
	c.Conn.Write([]byte("Error: " + err.Error() + "\r\n"))
}

func (c *client) Msg(msg string) {
	c.Conn.Write([]byte("> " + msg + "\r\n"))
}

/*
+------------------------------------------------------------------------------..+
[Chat rooms]-------|-[chat 1]---|---------|---------|---------|---------|---------|
Chat rooms:        | [User 1]: [Message_sdfnjasdfjajfdsjajfsdjlfkaslf]
*[chat 1]          | [User 2]: [Message_sdfnjasdfjaывфжадлфывжалжфывдлажфыжвалfsdf
[chat 2]           |    d;gsdgfjlkdfgjsdfjfdsjajfsdjlfkaslf]
[chat__long_name..]| [User 1]: [Message]
                   |
                   |

















*/
