package protocol

import "net"

type Room struct {
	name    string
	members map[net.Addr]Receiver
}

func (r *Room) Broadcast(sender Receiver, msg string) {
	for addr, member := range r.members {
		// if addr != sender.conn.RemoteAddr() {
		// 	member.msg(msg)
		// }
		if addr != sender.Conn().RemoteAddr() {
			//member.msg(msg)
			ReceiveMsg(member, msg)
		}
	}
}

func (r *Room) Join(member Receiver) {
	addr := member.Conn().RemoteAddr()
	r.members[addr] = member
}
