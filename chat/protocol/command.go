package protocol

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type commandID int

const (
	CMD_NICK = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type Command struct {
	ID      int
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

func Assemble(b []byte) (Command, error) {
	cmd_str := string(b)
	args := strings.Fields(cmd_str)
	if len(args) < 4 {
		return Command{}, fmt.Errorf("cann't assemble command: len(args) < 4")
	}
	cmd := Command{}
	cmd_id, err := strconv.Atoi(args[0])
	if err != nil {
		return Command{}, fmt.Errorf("cann't assemble command: cmd_id invalid (%v)", args[0])
	}
	cmd.ID = cmd_id

	adr, err := net.Dial(args[1], args[2])
	if err != nil {
		return Command{}, fmt.Errorf("cann't assemble command: net.Addr invalid: '%v' '%v'", args[1], args[2])
	}
	cmd.network = adr.RemoteAddr().Network()
	cmd.address = adr.RemoteAddr().String()

}
