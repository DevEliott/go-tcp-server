package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		make(map[string]*room),
		make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICKNAME:
			s.nickname(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		id:       uuid.New(),
		nickname: "anonymous",
		room:     nil,
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) nickname(c *client, args []string) {
	c.nickname = args[1]
	c.msg("nickname set to " + c.nickname)
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[uuid.UUID]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.id] = c

	s.quitRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nickname))
	c.msg("Welcome to " + roomName)
}

func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg("Available rooms are: " + strings.Join(rooms, ", "))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join a room first"))
		return
	}
	if len(args) < 2 {
		c.err(errors.New("missing message"))
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, fmt.Sprintf("%s: %s", c.nickname, msg))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("%s client %s disconnected", c.conn.RemoteAddr(), c.id)

	s.quitRoom(c)

	c.msg("Goodbye!")

	c.conn.Close()
}

func (s *server) quitRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.id)
		c.room.broadcast(c, fmt.Sprintf("%s left the room", c.nickname))
	}
}
