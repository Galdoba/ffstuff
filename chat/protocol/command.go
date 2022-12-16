package protocol

import "github.com/ffstuff/chat/client"

type CommandID int

const (
	CMD_NICK CommandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type Command struct {
	ID     CommandID
	Client *client.Client
	Args   []string
}
