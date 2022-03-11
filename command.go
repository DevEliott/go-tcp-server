package main

type commandID int

const (
	CMD_NICKNAME commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_UNDEFINED
)

type command struct {
	id     commandID
	client *client
	args   []string
}
