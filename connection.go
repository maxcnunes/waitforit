package main

import (
	"regexp"
	"strconv"
)

const regexAddressConn string = `^([a-z]{3,}):\/\/([a-zA-Z0-9\.\-_]+):?([0-9]*)*([a-zA-Z0-9\/\.\-_\(\)?=&#%]*)*$`

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
		return nil
	}

	match := regexp.MustCompile(regexAddressConn).FindAllStringSubmatch(address, -1)
	if len(match) < 1 {
		return nil
	}

	res := match[0]
	conn := &Connection{
		Type: res[1],
		Host: res[2],
		Path: res[4],
	}

	if conn.Type != "tcp" {
		conn.Scheme = conn.Type
		conn.Type = "tcp"
	}

	// resolve port
	if len(res[3]) == 0 {
		if conn.Scheme == "https" {
			conn.Port = 443
		} else {
			conn.Port = 80
		}
	} else {
		if port, err := strconv.Atoi(res[3]); err == nil {
			conn.Port = port
		}
	}

	return conn
}
