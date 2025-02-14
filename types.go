package main

import (
	"github.com/andybzn/gator/internal/config"
	"github.com/charmbracelet/log"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command, *log.Logger) error
}
