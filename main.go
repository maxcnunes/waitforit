package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var debug *bool

const pattConn string = `^([a-z]{3,}):\/\/([^:]+):?([0-9]+)?$`

func logDebug(msg interface{}) {
	if *debug {
		log.Print(msg)
	}
}

// Connection data
type Connection struct {
	Type   string
	Scheme string
	Port   int
	Host   string
}

func buildConn(host string, port int, fullConn string) *Connection {
	if host != "" {
		return &Connection{Type: "tcp", Host: host, Port: port}
	}

	if fullConn == "" {
		return nil
	}

	res := regexp.MustCompile(pattConn).FindAllStringSubmatch(fullConn, -1)[0]
	if len(res) != 4 {
		return nil
	}

	port, err := strconv.Atoi(res[3])
	if err != nil {
		port = 80
	}

	conn := &Connection{Type: res[1], Host: res[2], Port: port}
	if conn.Type != "tcp" {
		conn.Scheme = conn.Type
		conn.Type = "tcp"
	}

	if conn.Scheme == "https" {
		conn.Port = 443
	}

	return conn
}

func pingTCP(conn *Connection, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	logDebug("Dial adress: " + address)

	for {
		_, err := net.DialTimeout(conn.Type, address, time.Second)
		logDebug("ping TCP")

		if err == nil {
			logDebug("Up")
			return nil
		}

		logDebug("Down")
		logDebug(err)
		if time.Since(start) > timeout {
			return err
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func pingHTTP(conn *Connection, timeoutSeconds int) error {
	timeout := time.Duration(timeoutSeconds) * time.Second
	start := time.Now()
	address := fmt.Sprintf("%s://%s:%d", conn.Scheme, conn.Host, conn.Port)
	logDebug("HTTP adress: " + address)

	for {
		resp, err := http.Get(address)
		logDebug("ping HTTP " + resp.Status)

		if err == nil && resp.StatusCode < http.StatusInternalServerError {
			return nil
		}

		if time.Since(start) > timeout {
			return errors.New(resp.Status)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	fullConn := flag.String("full-connection", "", "full connection")
	host := flag.String("host", "", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "time to wait until port become available")
	debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	conn := buildConn(*host, *port, *fullConn)
	if conn == nil {
		log.Fatal("Invalid connection")
	}

	logDebug("Waiting " + strconv.Itoa(*timeout) + " seconds")
	if err := pingTCP(conn, *timeout); err != nil {
		log.Fatal(err)
	}

	if conn.Scheme != "http" && conn.Scheme != "https" {
		return
	}

	if err := pingHTTP(conn, *timeout); err != nil {
		log.Fatal(err)
	}
}
