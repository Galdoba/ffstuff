package protocol

import (
	"fmt"
	"net"
)

type Room struct {
	name    string
	members map[net.Addr]Receiver
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		members: make(map[net.Addr]Receiver),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Members() map[net.Addr]Receiver {
	return r.members
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

func (r *Room) Join(member Receiver) error {
	addr := member.Conn().RemoteAddr()
	if _, ok := r.members[addr]; ok == true {
		return fmt.Errorf("cann't join room: same connection already exists")
	}
	r.members[addr] = member
	return nil
}

func (r *Room) Remove(member Receiver) error {
	addr := member.Conn().RemoteAddr()
	if _, ok := r.members[addr]; ok == false {
		return fmt.Errorf("cann't remove from room: connection not mapped")
	}
	//r.members[addr] = member
	delete(r.members, addr)
	return nil
}
