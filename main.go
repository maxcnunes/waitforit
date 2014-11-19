package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func Dial(host string, port int, timeout int) error {
	_timeout := time.Duration(timeout) * time.Second

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), _timeout)
	if err != nil {
		return err
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	_, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}

	return nil
}

func main() {
	host := flag.String("host", "localshot", "host to connect")
	port := flag.Int("port", 80, "port to connect")
	timeout := flag.Int("timeout", 10, "timeout to wait port be available")

	flag.Parse()

	fmt.Println("starting")
	if err := Dial(*host, *port, *timeout); err != nil {
		// fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
