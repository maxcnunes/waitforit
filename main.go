package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

var debug *bool

const pattConn string = `([a-z]{3}):\/\/(.+):([0-9]+)`

func logDebug(msg string, conn *Connection) {
	if *debug {
		log.Printf("%s - %s://%s", msg, conn.Type, conn.Address)
	}
}

type Connection struct {
	Type    string
	Address string
}

func buildConn(host string, port int, fullConn string) *Connection {
	if host != "" {
		return &Connection{Type: "tcp", Address: fmt.Sprintf("%s:%d", host, port)}
	}

	if fullConn == "" {
		return nil
	}

	res := regexp.MustCompile(pattConn).FindAllStringSubmatch(fullConn, -1)[0]
	if len(res) != 4 {
		return nil
	}

	return &Connection{Type: res[1], Address: fmt.Sprintf("%s:%s", res[2], res[3])}
}

func dial(conn *Connection, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()

	for {
		_, err := net.Dial(conn.Type, conn.Address)
		if err == nil {
			logDebug("Up", conn)
			return nil
		}

		logDebug("Down", conn)
		if time.Since(start) > timeout {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func main() {
	fullConn := flag.String("full-connection", "", "full connection")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 0, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until port become available")
	debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	conn := buildConn(*host, *port, *fullConn)
	if conn == nil {
		log.Fatal("Invalid connection")
	}

	logDebug("Waiting", conn)
	if err := dial(conn, *timeout); err != nil {
		log.Fatal(err)
	}
}
