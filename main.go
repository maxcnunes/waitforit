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

const regexAddressConn string = `^([a-z]{3,}):\/\/([^:]+):?([0-9]+)?$`
const regexPathAddressConn string = `^([^\/]+)(\/?.*)$`

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
	logDebug("Dial address: " + address)

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
	address := fmt.Sprintf("%s://%s:%d%s", conn.Scheme, conn.Host, conn.Port, conn.Path)
	logDebug("HTTP address: " + address)

	for {
		resp, err := http.Get(address)

		if resp != nil {
			logDebug("ping HTTP " + resp.Status)
		}

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
	version := flag.Bool("v", false, "show the current version")
	debug = flag.Bool("debug", false, "enable debug")

	flag.Parse()

	if *version {
		fmt.Println("waitforit version 1.3.1")
		return
	}

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
