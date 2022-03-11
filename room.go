package main

import (
	"github.com/google/uuid"
)

type room struct {
	name    string
	members map[uuid.UUID]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if addr != sender.id {
			m.msg(msg)
		}
	}
}
