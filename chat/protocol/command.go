package protocol

import (
	"fmt"
	"net"
	"strings"
)

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_ERR
)

type Command struct {
	ID      commandID
	Sender  Receiver
	network string
	address string
	Args    []string
}

type Receiver interface {
	Conn() net.Conn
}

func ReceiveMsg(r Receiver, msg string) {
	r.Conn().Write([]byte("> " + msg + "\n"))
}

func Assemble(cmd_str string) (Command, error) {
	args := strings.Fields(cmd_str)
	if len(args) < 4 {
		return Command{}, fmt.Errorf("cann't assemble command: len(args) < 4")
	}
	cmd := Command{}
	cmd_id := convertCommandID(args[0])
	if cmd_id == CMD_ERR {
		return Command{}, fmt.Errorf("cann't assemble command: cmd_id invalid (%v)", args[0])
	}
	cmd.ID = cmd_id
	//adr, err := net.Dial(args[1], args[2])
	//if err != nil {
	//	fmt.Println("cann't assemble command: net.Addr invalid: ", args[1], args[2])
	//		return nil, fmt.Errorf("cann't assemble command: net.Addr invalid: '%v' '%v'", args[1], args[2])
	//	}

	cmd.network = args[1]
	cmd.address = args[2]
	cmd.Args = args[3:]
	return cmd, nil
}

func convertCommandID(str string) commandID {
	switch str {
	case "/nick":
		return CMD_NICK
	case "/join":
		return CMD_JOIN
	case "/rooms":
		return CMD_ROOMS
	case "/msg":
		return CMD_MSG
	case "/quit":
		return CMD_QUIT
	default:
		return CMD_ERR
	}
	/*
	   CMD_NICK
	   CMD_JOIN
	   CMD_ROOMS
	   CMD_MSG
	   CMD_QUIT
	*/
}
