package protocol

import (
	"net"

	"github.com/ffstuff/chat/client"
)

type room struct {
	Name    string
	Members map[net.Addr]*client.Client
}

func (r *room) Broadcast(sender *client.Client, msg string) {
	for addr, member := range r.Members {
		if addr != sender.conn.RemoteAddr() {
			member.Msg(msg)
		}
	}
}
