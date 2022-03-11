package main

import (
	"bufio"
	"errors"
	"net"
	"strings"

	"github.com/google/uuid"
)

type client struct {
	conn     net.Conn
	id       uuid.UUID
	nickname string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		var id commandID
		switch cmd {
		case "/nickname":
			id = CMD_NICKNAME
		case "/join":
			id = CMD_JOIN
		case "/rooms":
			id = CMD_ROOMS
		case "/msg":
			id = CMD_MSG
		case "/quit":
			id = CMD_QUIT
		default:
			id = CMD_UNDEFINED
			c.err(errors.New("unknown command" + cmd))
		}
		if id != CMD_UNDEFINED {
			c.commands <- command{
				id:     id,
				client: c,
				args:   args,
			}
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
