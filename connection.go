package main

import (
	"regexp"
	"strconv"
)

const regexAddressConn string = `^([a-z]{3,}):\/\/([^:]+):?([0-9]+)?$`
const regexPathAddressConn string = `^([^\/]+)(\/?.*)$`

// Connection data
type Connection struct {
	Type   string
	Scheme string
	Port   int
	Host   string
	Path   string
}

// BuildConn build a connection structure.
// This connection data can later be used as a common structure
// by the functions that will check if the target is available.
func BuildConn(cfg *Config) *Connection {
	if cfg.Host != "" {
		return &Connection{Type: "tcp", Host: cfg.Host, Port: cfg.Port}
	}

	address := cfg.Address
	if address == "" {
		// compatibilty with old argument
		address = cfg.FullConn
	}
	if address == "" {
		return nil
	}

	match := regexp.MustCompile(regexAddressConn).FindAllStringSubmatch(address, -1)
	if len(match) < 1 {
		return nil
	}

	res := match[0]

	port, err := strconv.Atoi(res[3])
	if err != nil {
		port = 80
	}

	hostAndPath := regexp.MustCompile(regexPathAddressConn).FindAllStringSubmatch(res[2], -1)[0]
	conn := &Connection{
		Type: res[1],
		Port: port,
		Host: hostAndPath[1],
		Path: hostAndPath[2],
	}

	if conn.Type != "tcp" {
		conn.Scheme = conn.Type
		conn.Type = "tcp"
	}

	if conn.Scheme == "https" {
		conn.Port = 443
	}

	return conn
}
