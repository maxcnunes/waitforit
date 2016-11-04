package main

import (
	"regexp"
	"strconv"
)

const regexAddressConn string = `^([a-z]{3,}):\/\/([^:]+):?([0-9]+)?$`
const regexPathAddressConn string = `^([^\/]+)(\/?.*)$`
const tcp = "tcp"
const https = "https"

// Connection data
type Connection struct {
	Type   string
	Scheme string
	Port   int
	Host   string
	Path   string
}

func buildConn(host string, port int, fullConn string) *Connection {
	if host != "" {
		return &Connection{Type: "tcp", Host: host, Port: port}
	}

	if fullConn == "" {
		return nil
	}

	res := regexp.MustCompile(regexAddressConn).FindAllStringSubmatch(fullConn, -1)[0]
	if len(res) != 4 {
		return nil
	}

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

	if conn.Type != tcp {
		conn.Scheme = conn.Type
		conn.Type = tcp
	}

	if conn.Scheme == https {
		conn.Port = 443
	}

	return conn
}
